package zerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	CodeUndefined = `zerror:undefined`
	CodeInternal  = `zerror:internal`
)

var (
	registered int32
	manager    *Manager
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
	HttpCode    int          `json:"http_code"`
	Msg         string       `json:"-"`
	LogLevel    logrus.Level `json:"-"`
	Description string       `json:"description"`
}

type Manager struct {
	Options
	errGroups []interface{}
}

type Responser interface {
	SetCode(code string)
	SetMessage(msg string)
}

type StdResponse struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  *string     `json:"msg"`
}

func (r *StdResponse) SetCode(code string) {
	r.Code = code
}

func (r *StdResponse) SetMessage(msg string) {
	r.Msg = &msg
}

type zerror struct {
	callLocation string
	callerName   string
	cause        error
	def          *Def
}


type withMessage struct {
	cause error
	msg   string
}


func (ze *zerror) Cause() error { return ze.cause }
func (ze *zerror) Error() string {
	if ze.cause == nil {
		return ze.def.Msg
	}
	return ze.def.Msg + ": " + ze.cause.Error()
}



func (ze *zerror) GetCaller() (string, string) {
	return ze.callLocation, ze.callerName
}


func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}


func DefaultDef(msg string) *Def {
	return &Def{
		Code:        "",
		HttpCode:    manager.defaultHttpCode,
		Msg:         msg,
		LogLevel:    manager.defaultLogLevel,
		Description: "",
	}
}


func (def *Def) SetCode(code string) *Def {
	def.Code = code
	return def
}

func (def *Def) SetHttpCode(code int) *Def {
	def.HttpCode = code
	return def
}

func (def *Def) SetMsg(msg string) *Def {
	def.Msg = msg
	return def
}

func (def *Def) SetLogLevel(level logrus.Level) *Def {
	def.LogLevel = level
	return def
}

func (def *Def) SetDesc(desc string) *Def {
	def.Description = desc
	return def
}

func (def *Def) wrap(err error, skip int) *zerror {

	l, n := getCaller(def, skip)
	org, ok := err.(*zerror)
	if ok {
		zerr := &zerror{
			callLocation: org.callLocation,
			callerName:   n + `/` + org.callerName,
			cause:        org,
			def:          def,
		}
		// if the original error is internal ,then the final error is internal
		if org.def.Code == CodeInternal {
			zerr.def = org.def
		}
		return zerr
	}
	return &zerror{
		callLocation: l,
		callerName:   n,
		cause:        err,
		def:          def,
	}
}

func (def *Def) Wrap(err error) *zerror {
	return def.wrap(err, 3)
}

func (def *Def)Wrapf(err error, format string, args ...interface{}) *zerror {
	wErr := &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return def.wrap(wErr, 3)
}

func (def *Def) New(msg string) *zerror {
	err := errors.New(msg)
	return def.wrap(err, 3)
}

func (def *Def)Errorf(format string, args ...interface{})*zerror {
	err := errors.New(fmt.Sprintf(format, args...))
	return def.wrap(err, 3)
}

// if err is not zerr.Def, then this err will be used
var UndefinedError = &Def{
	Code:        CodeUndefined,
	HttpCode:    500,
	Msg:         "unkown error",
	LogLevel:    logrus.ErrorLevel,
	Description: "error not defined, please contact admin",
}

var InternalError = &Def{
	Code:        CodeInternal,
	HttpCode:    500,
	Msg:         `internal error`,
	LogLevel:    logrus.ErrorLevel,
	Description: `this is server internal error, please contact admin`,
}


func (def *Def) Log(err error){
	def.log(err, 2)
}

//
func JSON(c *gin.Context, err error) {
	if registered == 0 {
		panic(`groups not registered`)
	}
	var def *Def
	zerr, ok := err.(*zerror)
	if !ok {
		def = UndefinedError
		zerr = &zerror{
			callLocation: "",
			callerName:   "",
			cause:        err,
			def:          def,
		}
		zerr.callLocation, zerr.callerName = getCaller(def, 2)
	} else {
		def = zerr.def
	}

	c.JSON(def.HttpCode, def.GetResponser())
	c.Abort()
	l, n := zerr.GetCaller()
	fields := logrus.Fields{`caller`: n}
	if l != `` {
		fields[`call_location`] = l
	}

	manager.logger.WithFields(fields).WithError(zerr.cause).Log(def.LogLevel, def.Msg)
}


func (def *Def) JSON(c *gin.Context, err error) {
	if registered == 0 {
		panic(`groups not registered`)
	}

	httpCode := def.HttpCode
	c.JSON(httpCode, def.GetResponser())
	c.Abort()
	def.log(err, 3)
}


func (def *Def)GetResponser() Responser {
	s := manager.responseFunc()
	s.SetCode(def.Code)
	if manager.RespondMessage {
		s.SetMessage(def.Description)
	}
	return s
}

func (def *Def) log(err error, skip int){

	fields := logrus.Fields{}
	l, n := getCaller(def, skip)
	fields[`caller`] = n
	if l != `` {
		fields[`call_location`] = l
	}
	manager.logger.WithFields(fields).WithError(err).Log(def.LogLevel, def.Msg)
}

func getCaller(def *Def, skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)
	var callLocation, callerName string
	if ok && (manager.debug && def.LogLevel == logrus.DebugLevel || def.Code == CodeInternal || def.Code == CodeUndefined) {
		callLocation = file + "/" + strconv.Itoa(line)
	}
	if ok {
		funcNameSplited := strings.Split(runtime.FuncForPC(pc).Name(), `.`)
		callerName = funcNameSplited[len(funcNameSplited)-1]
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

// the parameters must be error group ptr,
// if error group has field `Prefix`, then it's values will be used as error code prefix,
// else the group Type Name (after standardized) will be used as prefix,
// the prefix and the suberrorcode will be joined by ':'

func InitErrGroup(group interface{}) {
	typ := reflect.TypeOf(group)
	val := reflect.ValueOf(group)
	if typ.Kind() != reflect.Ptr {
		logrus.Panicf(`moduleErr is not ptr but: %s`, typ.Kind())
	}
	typ = typ.Elem()
	val = val.Elem()
	if typ.Kind() != reflect.Struct {
		logrus.Panicf(`moduleErr not struct, but: %s`, typ.Kind())
	}
	prefix := GetStandardName(typ.Name()) + manager.codeConnector
	nameField, ok := typ.FieldByName(`Prefix`)
	if ok {
		if nameField.Type.Kind() != reflect.String {
			logrus.Panicf(`Name field not string type`)
		}
		nameVal := val.FieldByName(`Prefix`).Interface().(string)
		if nameVal == `` {
			prefix = ``
		} else {
			prefix = nameVal + manager.codeConnector
		}
	}
	var zerr *Def
	for i := 0; i < typ.NumField(); i++ {
		tField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		if tField.Type != reflect.TypeOf(zerr) {
			continue
		}
		if structField.IsNil() {
			logrus.Panicf(`%s is nil`, tField.Name)
		}
		zerr = structField.Interface().(*Def)

		if zerr.Code != `` {
			continue
		}
		zerr.Code = fmt.Sprintf(`%s%s`, prefix, GetStandardName(tField.Name))
	}
}

func JsonDumpGroups(ident string) string {
	mared, err := json.MarshalIndent(manager.errGroups, ``, ident)
	if err != nil {
		panic(err)
	}
	return string(mared)
}

func New(options ...Option) *Manager {
	do := &Options{
		wordConnector:  `-`,
		codeConnector:  `:`,
		RespondMessage: true,
		logger:         logrus.StandardLogger(),
		responseFunc: func() Responser {
			return new(StdResponse)
		},
		defaultLogLevel: logrus.ErrorLevel,
		defaultHttpCode: 400,
	}
	for _, setter := range options {
		setter(do)
	}
	m := &Manager{
		Options: *do,
	}
	manager = m
	return m
}

func (m *Manager) RegisterGroups(groups ...interface{}) {
	if !atomic.CompareAndSwapInt32(&registered, 0, 1) {
		panic(`groups registered twice`)
	}
	for _, v := range groups {
		InitErrGroup(v)
		m.errGroups = append(m.errGroups, v)
	}
}

// for test
func unregister() {
	registered = 0
	manager.errGroups = nil
}
