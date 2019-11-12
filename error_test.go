package zerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestErr struct {
	TestErr1            *Def
	Err                 *Def
	ThisISAVeryLongName *Def
}

type TestErr1 struct {
	Err    *Def
	Prefix string
}

func TestMain(m *testing.M) {
	manager = New()
	m.Run()
}

func TestGenerateCode(t *testing.T) {
	data := &TestErr{
		TestErr1:            new(Def),
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`},
	}
	InitErrGroup(data)
	require.Equal(t, `test-err:test-err1`, data.TestErr1.Code)
	require.Equal(t, `test-err:this-is-avery-long-name`, data.ThisISAVeryLongName.Code)

	require.Equal(t, `custom-code`, data.Err.Code)
	data.Err.Code = ``
	InitErrGroup(data)
	require.Equal(t, `test-err:err`, data.Err.Code)

	data1 := &TestErr1{
		Err:    new(Def),
		Prefix: "",
	}
	InitErrGroup(data1)
	require.Equal(t, `err`, data1.Err.Code)

	data1.Prefix = `custom-prefix`
	data1.Err.Code = ``
	InitErrGroup(data1)
	require.Equal(t, `custom-prefix:err`, data1.Err.Code)

}

func ExampleJsonDumpGroups() {

	data := &TestErr{
		TestErr1:            new(Def),
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`},
	}
	manager.RegisterGroups(data)
	fmt.Println(JsonDumpGroups(``))
	manager.errGroups = nil
	// Output:
	// [
	// {
	// "TestErr1": {
	// "code": "test-err:test-err1",
	// "http_code": 0,
	// "description": ""
	// },
	// "Err": {
	// "code": "custom-code",
	// "http_code": 0,
	// "description": ""
	// },
	// "ThisISAVeryLongName": {
	// "code": "test-err:this-is-avery-long-name",
	// "http_code": 0,
	// "description": ""
	// }
	// }
	// ]
}

func ExampleGetCaller() {
	_, caller := getCaller(InternalError, 1)
	fmt.Println(caller)
	// Output:
	// ExampleGetCaller
}

func ExampleNested() {
	data := &TestErr{
		TestErr1:            &Def{Msg: `msg1`},
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`, Msg: `msg2`},
	}
	InitErrGroup(data)
	ze := data.TestErr1.Wrap(data.Err.Wrap(errors.New(`original-error`)))
	fmt.Println(ze.Error(), ze.callerName, ze.def.Code)
	// Output:
	// msg1: msg2: original-error ExampleNested/ExampleNested test-err:test-err1
}

type customeRsp struct {
	A   string
	Msg string
}

func (c *customeRsp) SetCode(code string) {
	c.A = code
}

func (c *customeRsp) SetMessage(msg string) {
	c.Msg = msg
}

func ExampleCustomResponser() {
	unregister()
	m := New(
		WithResponser(func() Responser {
			return new(customeRsp)
		}),
	)
	m.RegisterGroups()
	defer unregister()
	rsp := InternalError.GetResponser()
	mared, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(mared))
	// Output:
	// {"A":"zerror:internal","Msg":"this is server internal error, please contact admin"}
}

func ExampleDefaultDef(){

	unregister()
	m := New(DefaultHttpCode(500), DefaultLogLevel(logrus.InfoLevel))
	errDef := DefaultDef(`msg`)
	m.RegisterGroups()
	defer unregister()
	fmt.Printf("%+v\n", errDef)
	// Output:
	// &{Code: HttpCode:500 Msg:msg LogLevel:info Description:}
}

func ExampleLog(){
	unregister()
	logger := logrus.New()
	format := logrus.TextFormatter{
		DisableTimestamp:          true,
	}
	logger.SetFormatter(&format)
	m := New(Logger(logger))
	m.RegisterGroups()
	defer unregister()
	originalError := errors.New(`original error`)
	def := DefaultDef(`default`)
	def.Log(originalError)

	fmt.Println(def.New(`new`))

	errorf := def.Errorf(`%s`, `errorf`)
	fmt.Println(errorf.callerName, errorf.Error(), errorf.def.Msg)
	// Output:
	// default: new
	// ExampleLog default: errorf default
}
