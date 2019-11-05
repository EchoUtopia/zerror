package zerror

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

type stdResponse struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  *string     `json:"msg"`
}

func (r *stdResponse) SetCode(code string) {
	r.Code = code
}

func (r *stdResponse) SetMessage(msg string) {
	r.Msg = &msg
}

type zerror struct {
	callLocation string
	callerName   string
	cause        error
	def          *Def
	internal     bool
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

func (def *Def) Wrap(err error) *zerror {
	l, n := getCaller(def.LogLevel, 2)
	org, ok := err.(*zerror)
	if ok {
		return &zerror{
			callLocation: org.callLocation,
			callerName:   n + `/` + org.callerName,
			cause:        org,
			def:          def,
			internal:     org.internal,
		}
	}
	return &zerror{
		callLocation: l,
		callerName:   n,
		cause:        err,
		def:          def,
	}
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
			internal:     true,
		}
		zerr.callLocation, zerr.callerName = getCaller(def.LogLevel, 2)
	} else {
		def = zerr.def
	}

	c.JSON(def.HttpCode, getResponse(def))
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
	c.JSON(httpCode, getResponse(def))
	c.Abort()
	fields := logrus.Fields{}
	l, n := getCaller(def.LogLevel, 2)
	fields[`caller`] = n
	if def.LogLevel == logrus.DebugLevel {
		fields[`call_location`] = l
	}
	manager.logger.WithFields(fields).WithError(err).Log(def.LogLevel, def.Msg)
}

func getResponse(def *Def) Responser {
	s := manager.responseFunc()
	s.SetCode(def.Code)
	if manager.RespondMessage {
		s.SetMessage(def.Description)
	}
	return s
}

func getCaller(debugLevel logrus.Level, skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)
	var callLocation, callerName string
	if ok && debugLevel == logrus.DebugLevel {
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
			return new(stdResponse)
		},
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
