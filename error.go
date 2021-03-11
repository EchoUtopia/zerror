package zerror

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	ExtLogLvl = `log_level`
	ExtLogger = `logger`
)

type Def struct {
	// error code,
	Code        string       `json:"code"`
	Msg         string       `json:"msg"`
	Description string       `json:"description"`
	PCode       ProtocolCode `json:"protocol_code"`

	// extended fields
	extensions map[string]interface{}
}

type Data map[string]interface{}

type ZContext struct {
	callLocation string
	callerName   string
	Data         Data
	Ctx          context.Context
}

type Error struct {
	cause error
	*Def
	msg string
	*ZContext
}

func (ctx *ZContext) Merge(m *ZContext) {
	ctx.callLocation = m.callLocation
	ctx.callerName += `/` + m.callerName
	for k, v := range m.Data {
		ctx.Data[k] = v
	}
}

func (ze *Error) Unwrap() error {
	return ze.cause
}
func (ze *Error) Error() string {
	msg := ze.Def.Code

	if ze.msg != `` {
		msg += `: ` + ze.msg
	}

	if ze.cause != nil {
		msg += ` | ` + ze.cause.Error()
	}

	return msg
}

func (ze *Error) WithData(data Data) *Error {
	if ze.Data == nil {
		ze.Data = make(Data, len(data)+2)
	}
	for k, v := range data {
		ze.Data[k] = v
	}
	return ze
}

func (ze *Error) WithContext(ctx context.Context) *Error {
	ze.Ctx = ctx
	return ze
}

func (ze *Error) GetCaller() (string, string) {
	return ze.callLocation, ze.callerName
}

func (ze *Error) Render() Render {
	def := ze.Def
	s := renderPool.Get().(Render)
	code := def.Code
	if code != Internal.Code && Internal.Cause(ze.cause) {
		code = Internal.Code
	}
	s.SetCode(code)
	if Manager.RespondMsgSet && Manager.RespondMessage ||
		!Manager.RespondMsgSet && Manager.debugMode {
		s.SetMessage(ze.Error())
	}
	return s
}

func (def *Def) Extend(k string, v interface{}) *Def {
	if def.extensions == nil {
		def.extensions = make(map[string]interface{})
	}
	def.extensions[k] = v
	return def
}

func (def *Def) GetExtension(key string) (interface{}, bool) {
	if def.extensions == nil {
		return nil, false
	}
	value, ok := def.extensions[key]
	return value, ok
}

func (def *Def) wrapf(err error, skip int, format string, args ...interface{}) *Error {

	l, n := GetCaller(def, skip)
	zCause := &Error{}
	zErr := &Error{
		cause: err,
		Def:   def,
	}
	if ok := errors.As(err, &zCause); ok {
		zErr.ZContext = zCause.ZContext
		zErr.callerName += `/` + n
	} else {
		zErr.ZContext = &ZContext{
			callLocation: l,
			callerName:   n,
			Data:         make(Data),
		}
	}
	if format != `` {
		zErr.msg = fmt.Sprintf(format, args...)
	}
	return zErr
}

func (def *Def) Wrap(err error) *Error {
	return def.wrapf(err, 3, ``)
}

func (def *Def) Wrapf(err error, format string, args ...interface{}) *Error {
	return def.wrapf(err, 3, format, args...)
}

func (def *Def) WithMsg(msg string) *Error {
	return def.wrapf(nil, 3, msg)
}

func (def *Def) New() *Error {
	return def.wrapf(nil, 3, ``)
}

func (def *Def) Errorf(format string, args ...interface{}) *Error {
	err := errors.New(fmt.Sprintf(format, args...))
	return def.wrapf(err, 3, ``)
}

var renderPool = sync.Pool{
	New: func() interface{} {
		render := Manager.render()
		reset, ok := render.(Resetter)
		if ok {
			reset.Reset()
		}
		return render
	},
}

func (def *Def) Cause(err error) bool {
	zerr := &Error{}
	for {
		if ok := errors.As(err, &zerr); ok {
			if zerr.Def == def {
				return true
			} else {
				err = zerr.cause
			}
		} else {
			return false
		}
	}
}

func FromCode(code string) (bool, *Error) {
	def, ok := defMap[code]
	if !ok {
		return ok, nil
	}
	return true, def.New()
}

func GetCaller(def *Def, skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)
	var callLocation, callerName string
	if ok && (Manager.debugMode || def == nil || def.Code == BizCodeInternal) {
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
			out += Manager.wordConnector
			lastLower = false
		} else {
			lastLower = true
		}
		out += string(v)
	}
	lowered := strings.ToLower(out)
	return lowered
}
