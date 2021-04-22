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

type Def struct {
	// error code,
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	Description string `json:"desc"`
	// can be used as http status code or grpc status code
	Status Status `json:"status"`

	// extended fields
	extensions map[string]interface{}
}

type Data map[string]interface{}

type ZContext struct {
	callerLoc  string
	callerName string
	Data       Data
	Ctx        context.Context
}

func (ctx *ZContext) Merge(m *ZContext) {
	ctx.callerLoc = m.callerLoc
	ctx.callerName += `/` + m.callerName
	for k, v := range m.Data {
		ctx.Data[k] = v
	}
}

type Error struct {
	cause error
	*Def
	msg string
	ZContext
}

func (ze *Error) Unwrap() error {
	return ze.cause
}
func (ze *Error) Error() string {
	b := strings.Builder{}
	b.WriteString(ze.Def.Code)
	if ze.msg != `` {
		b.WriteString(`(`)
		b.WriteString(ze.msg)
		b.WriteString(`)`)
	}
	if ze.cause != nil {
		b.WriteString(` | `)
		b.WriteString(ze.cause.Error())
	}
	return b.String()
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

func (ze *Error) WithKVs(kvs ...interface{}) *Error {
	if ze.Data == nil {
		ze.Data = make(Data, len(kvs))
	}
	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			k = fmt.Sprintf(`%v`, kvs[i])
		}
		var v interface{}
		if i+1 < len(kvs) {
			v = kvs[i+1]
		}
		ze.Data[k] = v
	}
	return ze
}

func (ze *Error) WithCtx(ctx context.Context) *Error {
	ze.Ctx = ctx
	return ze
}

func (ze *Error) GetCaller() (string, string) {
	return ze.callerLoc, ze.callerName
}

func (ze *Error) Render() Render {
	def := ze.Def
	s := renderPool.Get().(Render)
	code := def.Code
	if code != Internal.Code && Internal.Cause(ze.cause) {
		code = Internal.Code
	}
	s.SetCode(code)
	if Manager.respondMsgSet && Manager.respondMessage ||
		!Manager.respondMsgSet && Manager.debugMode {
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
			callerLoc:  l,
			callerName: n,
			Data:       make(Data),
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

func FromCode(code string) (*Error, bool) {
	def, ok := defMap[code]
	if !ok {
		return nil, ok
	}
	return def.New(), true
}

func GetCaller(def *Def, skip int) (loc string, name string) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok && (Manager.debugMode || def == nil || def.Code == CodeInternal) {
		loc = file + "/" + strconv.Itoa(line)
	}
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		name = funcName[strings.LastIndexByte(funcName, '.')+1:]
	}
	return
}

// add '-' before initial capital letters and turn lower
func getStandardName(name string) string {
	out := ``
	lastLower := true
	for k, v := range name {
		if v >= 'A' && v <= 'Z' && k != 0 && lastLower {
			out += `-`
			lastLower = false
		} else {
			lastLower = true
		}
		out += string(v)
	}
	lowered := strings.ToLower(out)
	return lowered
}
