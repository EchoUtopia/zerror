package zerror

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync/atomic"
)

var (
	registered int32
	Manager    = &Zmanager{Options: &Options{defaultPCode: CodeInvalid}}
)

type Zmanager struct {
	*Options
	errGroups []interface{}
}

// the parameters must be error group ptr,
// if error group has field `Prefix`, then it's values will be used as error code prefix,
// else the group Type Name (after standardized) will be used as prefix,
// the prefix and the suberrorcode will be joined by ':'

func InitErrGroup(group interface{}) {
	typ := reflect.TypeOf(group)
	val := reflect.ValueOf(group)
	if typ.Kind() != reflect.Ptr {
		log.Panicf(`error group is not ptr but: %s`, typ.Kind())
	}
	typ = typ.Elem()
	val = val.Elem()
	if typ.Kind() != reflect.Struct {
		log.Panicf(`error group is not struct, but: %s`, typ.Kind())
	}

	groupName := typ.Name()
	if !strings.HasSuffix(typ.Name(), `Group`) || len(groupName) <= 5 {
		log.Panicf(`error group type: %s must has suffix Group`, groupName)
	}
	prefix := GetStandardName(groupName[:len(groupName)-5]) + Manager.codeConnector
	nameField, ok := typ.FieldByName(`Prefix`)
	if ok {
		if nameField.Type.Kind() != reflect.String {
			log.Panicf(`error group: %s, Prefix field not string type, but: %s`, groupName, nameField.Type.Kind())
		}
		nameVal := val.FieldByName(`Prefix`).Interface().(string)
		if nameVal == `` {
			prefix = ``
		} else {
			prefix = nameVal + Manager.codeConnector
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
			log.Panicf(`error group: %s, field:%s is nil`, groupName, tField.Name)
		}
		zerr = structField.Interface().(*Def)

		if zerr.Code != `` {
			continue
		}
		if zerr.PCode == -1 {
			zerr.PCode = Manager.defaultPCode
		}
		zerr.Code = fmt.Sprintf(`%s%s`, prefix, GetStandardName(tField.Name))
	}
}

func (m *Zmanager) JsonDumpGroups(ident string) string {
	mared, err := json.MarshalIndent(m.errGroups, ``, ident)
	if err != nil {
		panic(err)
	}
	return string(mared)
}

func (m *Zmanager) GetErrorGroups() []interface{} {
	return m.errGroups
}

func New(options ...Option) *Zmanager {
	do := &Options{
		wordConnector:  `-`,
		codeConnector:  `:`,
		RespondMessage: true,
		responseFunc: func() Responser {
			return new(StdResponse)
		},
		defaultPCode: 400,
	}
	for _, setter := range options {
		setter(do)
	}
	m := &Zmanager{
		Options: do,
	}
	Manager = m
	return m
}

func (m *Zmanager) RegisterGroups(groups ...interface{}) {
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
	Manager.errGroups = nil
	Manager.defaultPCode = CodeInvalid
}

func Registered() bool {
	return registered == 1
}
