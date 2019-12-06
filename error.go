package zerror

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	BizCodeInternal = `zerror:internal`

	ExtLogLvl = `log_level`
)

// error can be used alone, like:
// var SmsCode = &zerror.Def{`sms:code`, 500, `sms code`, logrus.ErrorLevel, ``}
// then you can use it to respond and log error:
// SmsCode.JSON(ginContext, error)

// it's better to use error group, error group classfy errors into different groups:
// 	Common = &CommonErr{
//		Prefix:   "",
//		Args:     &zerror.Def{``, 400, `args err`, logrus.DebugLevel, ``},
//		Internal: &zerror.Def{``, 500, `internal err`, logrus.ErrorLevel, ``},
//	}

// then you can use it to repond and log error:
// Common.Args.JSON(c, errors.Wrap(err, `msg`))

// the error code will be automaticlly generated if you use error group, the name will be standardized
// if there's a 'Prefix' field in Group, then it's value will be used as error code prefix,
// else the group type name will be used,
// but if the value is empty, then there's no prefix in error code
// take 'Common' group above as example:
// the Common.Args's code will be `args`, because the Prefix's values is ''
// if CommonErr has no Prefix field, then the Common.Args's error code will be 'common-err:args'
// if the Prefix field's values is 'common', the the Common.Args's error code will be 'common:args'

type Def struct {
	// error code,
	Code        string       `json:"code"`
	Msg         string       `json:"-"`
	Description string       `json:"description"`
	PCode       ProtocolCode `json:"protocol_code"`

	// extended fields
	Extensions map[string]interface{} `json:"extensions"`
}

type Data map[string]interface{}

type ZContext struct {
	CallLocation string
	CallerName   string
	Data         Data
	Ctx          context.Context
}

type Error struct {
	cause error
	Def   *Def
	msg   string
	*ZContext
}

func NewError(cause error, def *Def, msg string, ctx *ZContext) *Error {
	return &Error{
		cause:    cause,
		Def:      def,
		msg:      msg,
		ZContext: ctx,
	}
}

func (ze *Error) Cause() error { return ze.cause }
func (ze *Error) Error() string {
	msg := ze.Def.Msg

	if ze.msg != `` {
		msg += `: ` + ze.msg
	}

	if ze.cause != nil {
		msg += ` <- ` + ze.cause.Error()
	}

	return msg
}

func (ze *Error) WithData(kvs ...*kv) *Error {
	if ze.Data == nil {
		ze.Data = make(Data, len(kvs)+2)
	}
	for _, v := range kvs {
		ze.Data[v.K] = v.V
	}
	return ze
}

type kv struct {
	K string
	V interface{}
}

func KV(k string, v interface{}) *kv {
	return &kv{
		K: k,
		V: v,
	}
}

func (ze *Error) WithContext(ctx context.Context) *Error {
	ze.Ctx = ctx
	return ze
}

func (ze *Error) GetCaller() (string, string) {
	return ze.CallLocation, ze.CallerName
}

func DefaultDef(msg string) *Def {
	return &Def{
		Code:        "",
		PCode:       manager.defaultPCode,
		Msg:         msg,
		Description: "",
	}
}

func (def *Def) WithCode(code string) *Def {
	def.Code = code
	return def
}

func (def *Def) WithPCode(pCode ProtocolCode) *Def {
	def.PCode = pCode
	return def
}

func (def *Def) WithMsg(msg string) *Def {
	def.Msg = msg
	return def
}

func (def *Def) Extend(k string, v interface{}) *Def {
	if def.Extensions == nil {
		def.Extensions = make(map[string]interface{})
	}
	def.Extensions[k] = v
	return def
}

func (def *Def) WithDesc(desc string) *Def {
	def.Description = desc
	return def
}

func (def *Def) wrapf(err error, skip int, format string, args ...interface{}) *Error {

	l, n := GetCaller(def, skip)
	org, ok := err.(*Error)
	var zerr *Error
	if ok {
		zerr = &Error{
			ZContext: &ZContext{
				CallLocation: org.CallLocation,
				CallerName:   n + `/` + org.CallerName,
			},
			cause: org,
			Def:   def,
		}
		// if the original error is internal ,then the final error is internal
		if org.Def.Code == BizCodeInternal {
			zerr.Def = org.Def
		}
		return zerr
	} else {
		zerr = &Error{
			ZContext: &ZContext{
				CallLocation: l,
				CallerName:   n,
			},
			cause: err,
			Def:   def,
		}
	}
	if format != `` {
		zerr.msg = fmt.Sprintf(format, args...)
	}
	return zerr
}

func (def *Def) Wrap(err error) *Error {
	return def.wrapf(err, 3, ``)
}

func (def *Def) Wrapf(err error, format string, args ...interface{}) *Error {
	return def.wrapf(err, 3, format, args...)
}

func (def *Def) New(msg string) *Error {
	err := errors.New(msg)
	return def.wrapf(err, 3, ``)
}

func (def *Def) Errorf(format string, args ...interface{}) *Error {
	err := errors.New(fmt.Sprintf(format, args...))
	return def.wrapf(err, 3, ``)
}

var InternalError = &Def{
	Code:        BizCodeInternal,
	PCode:       CodeInternal,
	Msg:         `internal error`,
	Description: `this is server internal error, please contact admin`,
}

func (def *Def) GetResponser(err error) Responser {
	s := manager.responseFunc()
	s.SetCode(def.Code)
	if manager.RespondMessage || manager.debugMode {
		s.SetMessage(err.Error())
	}
	return s
}

func (def *Def) Equal(err error) bool {
	zerr, ok := err.(*Error)
	if !ok {
		return false
	}
	return zerr.Def == def
}

func GetCaller(def *Def, skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)
	var callLocation, callerName string
	if ok && (manager.debugMode || def.Code == BizCodeInternal) {
		callLocation = file + "/" + strconv.Itoa(line)
	}
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		callerName = funcName[strings.LastIndexByte(funcName, '.')+1:]
	}
	return callLocation, callerName
}

// add '-' before initial capital letters and turn lower
func GetStandardName(name string) string {
	out := ``
	lastLower := true
	for k, v := range name {
		if v >= 'A' && v <= 'Z' && k != 0 && lastLower {
			out += manager.wordConnector
			lastLower = false
		} else {
			lastLower = true
		}
		out += string(v)
	}
	lowered := strings.ToLower(out)
	return lowered
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
