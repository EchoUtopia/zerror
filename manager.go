package zerror

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	Manager = &Zmanager{Options: &Options{defaultPCode: CodeInvalid}}
	defMap  = map[string]*Def{}
)

type Zmanager struct {
	*Options
	errGroups []interface{}
	sync.Mutex
	registered int32
}

// the parameters must be error group ptr,
// if error group has field `Prefix`, then it's values will be used as error code prefix,
// else the group Type Name (after standardized) will be used as prefix,
// the prefix and the suberrorcode will be joined by ':'

func initErrGroup(group interface{}) {
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

	groupName := GetStandardName(typ.Name())
	//if !strings.HasSuffix(typ.Name(), `Group`) || len(groupName) <= 5 {
	//	log.Panicf(`error group type: %s must has suffix Group`, groupName)
	//}
	//prefix := GetStandardName(groupName[:len(groupName)-5]) + Manager.codeConnector
	prefix := groupName + Manager.codeConnector
	nameField, ok := typ.FieldByName(`Prefix`)
	if ok {
		if nameField.Type.Kind() != reflect.String {
			log.Panicf(`error group: %s, Prefix field is not string type, but: %s`,
				groupName, nameField.Type.Kind())
		}
		nameVal := val.FieldByName(`Prefix`).Interface().(string)
		if nameVal == `` {
			prefix = ``
		} else {
			prefix = nameVal + Manager.codeConnector
		}
	}
	errCnt := 0
	var def *Def
	for i := 0; i < typ.NumField(); i++ {
		tField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		if tField.Type != reflect.TypeOf(def) {
			continue
		}
		if structField.IsNil() {
			def = &Def{
				Code: fmt.Sprintf(`%s%s`, prefix, GetStandardName(tField.Name)),
			}
			structField.Set(reflect.ValueOf(def))
		} else {
			def = structField.Interface().(*Def)
		}

		if def.PCode == CodeInvalid {
			def.PCode = Manager.defaultPCode
		}

		if def.Code == `` {
			def.Code = fmt.Sprintf(`%s%s`, prefix, GetStandardName(tField.Name))
		}
		if defMap[def.Code] != nil {
			log.Panicf(`def code: %s duplicated`, def.Code)
		}
		errCnt++
		defMap[def.Code] = def
	}
	if errCnt == 0 {
		log.Panicf(`error def not found in group: %s`, group)
	}
}

func (m *Zmanager) GetErrorGroups() []interface{} {
	if !m.Registered() {
		panic(`not registered`)
	}
	return m.errGroups
}

func Init(options ...Option) *Zmanager {
	do := &Options{
		wordConnector:  `-`,
		codeConnector:  `:`,
		RespondMessage: true,
		render: func() Render {
			return new(StdResponse)
		},
		defaultPCode: 200,
	}
	for _, setter := range options {
		setter(do)
	}
	m := &Zmanager{
		Options: do,
	}
	Manager = m
	return Manager
}

func (m *Zmanager) RegisterGroups(groups ...interface{}) {
	if m.Registered() {
		panic(`registered twice`)
	}
	m.Lock()
	defer m.Unlock()
	if m.registered == 1 {
		panic(`groups registered twice`)
	}
	for _, v := range groups {
		initErrGroup(v)
		m.errGroups = append(m.errGroups, v)
	}
	atomic.StoreInt32(&m.registered, 1)
}

// for test
func unregister() {
	Manager.registered = 0
	Manager.errGroups = nil
	Manager.defaultPCode = CodeInvalid
	defMap = map[string]*Def{}
}

func (m *Zmanager) Registered() bool {
	return atomic.LoadInt32(&m.registered) == 1
}
