package zerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestErr struct {
	Group
	TestErr1            *Def
	Err                 *Def
	ThisISAVeryLongName *Def
}

type TestErr1 struct {
	Group
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
	require.Equal(t, `test-err:test-err1`, data.TestErr1.PCode)
	require.Equal(t, `test-err:this-is-avery-long-name`, data.ThisISAVeryLongName.PCode)

	require.Equal(t, `custom-code`, data.Err.PCode)
	data.Err.Code = ``
	InitErrGroup(data)
	require.Equal(t, `test-err:err`, data.Err.PCode)

	data1 := &TestErr1{
		Err:    new(Def),
		Prefix: "",
	}
	InitErrGroup(data1)
	require.Equal(t, `err`, data1.Err.PCode)

	data1.Prefix = `custom-prefix`
	data1.Err.Code = ``
	InitErrGroup(data1)
	require.Equal(t, `custom-prefix:err`, data1.Err.PCode)

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
	_, caller := GetCaller(InternalError, 1)
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
	fmt.Println(ze.Error(), ze.CallerName, ze.Def.PCode)

	ze = data.TestErr1.Wrapf(data.Err.Wrap(errors.New(`original-error`)), `wrap message: %s`, `ad`)
	fmt.Println(ze.Error(), ze.CallerName, ze.Def.PCode)
	// Output:
	// original-error ExampleNested/ExampleNested test-err:test-err1
	// wrap message: ad: original-error ExampleNested/ test-err:test-err1
}

type customeRsp struct {
	A   string
	Msg string
}

func (c *customeRsp) SetBizCode(code string) {
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

func ExampleDefaultDef() {

	unregister()
	m := New(DefaultCode(500), DefaultLogLevel(InfoLevel))
	errDef := DefaultDef(`msg`)
	m.RegisterGroups()
	defer unregister()
	fmt.Printf("%+v\n", errDef)
	// Output:
	// &{Code: Msg:msg Description: LogLevel:info PCode:500 Extended:map[]}
}

func ExampleLog() {
	unregister()
	// logger := logrus.New()
	// format := logrus.TextFormatter{
	// 	DisableTimestamp: true,
	// }
	// logger.SetFormatter(&format)
	m := New()
	m.RegisterGroups()
	defer unregister()
	originalError := errors.New(`original error`)
	def := DefaultDef(`default`)
	Log(def, originalError)

	fmt.Println(def.New(`new`))

	errorf := def.Errorf(`%s`, `errorf`)
	fmt.Println(errorf.CallerName, errorf.Error(), errorf.Def.Msg)
	// Output:
	// new
	// ExampleLog errorf default
}

type ForDefault struct {
	Group
	Err1 *Def
}

func ExampleDefaultDef2() {

	unregister()
	d := &ForDefault{Err1: DefaultDef(`default`)}
	fmt.Println(d.Err1.PCode, d.Err1.LogLevel)
	m := New(DefaultLogLevel(ErrorLevel), DefaultCode(500))
	m.RegisterGroups(d)
	defer unregister()

	fmt.Println(d.Err1.PCode, d.Err1.LogLevel)
	// Output:
	// -1 unknown
	// 500 error
}

func ExampleDef_Is() {

	originalError := errors.New(`original error`)
	def := DefaultDef(`default`)
	def1 := DefaultDef(`default2`)
	wrapped := def.Wrap(originalError)
	wrapped1 := def1.Wrap(wrapped)
	fmt.Println(def.Equal(wrapped))
	fmt.Println(def1.Equal(wrapped1))
	fmt.Println(def.Equal(wrapped1))
	fmt.Println(def1.Equal(wrapped))
	// Output:
	// true
	// true
	// false
	// false
}

func BenchmarkDef_Wrap(b *testing.B) {

	originalError := errors.New(`original error`)
	def := DefaultDef(`default`)
	for i := 0; i < b.N; i++ {
		def.wrapf(originalError, 1, ``)
	}
}
