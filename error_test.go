package zerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestErrGroup struct {
	TestErr1            *Def
	Err                 *Def
	ThisISAVeryLongName *Def
}

type TestErr1Group struct {
	Err    *Def
	Prefix string
}

func TestMain(m *testing.M) {
	Manager = New()
	m.Run()
}

func TestGenerateCode(t *testing.T) {
	data := &TestErrGroup{
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

	data1 := &TestErr1Group{
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

	data := &TestErrGroup{
		TestErr1:            new(Def),
		ThisISAVeryLongName: new(Def).Extend(ExtLogLvl, 23),
		Err:                 &Def{Code: `custom-code`},
	}
	Manager.RegisterGroups(data)
	fmt.Println(Manager.JsonDumpGroups(``))
	Manager.errGroups = nil
	// Output:
	// [
	// {
	// "TestErr1": {
	// "code": "test-err:test-err1",
	// "description": "",
	// "protocol_code": 0
	// },
	// "Err": {
	// "code": "custom-code",
	// "description": "",
	// "protocol_code": 0
	// },
	// "ThisISAVeryLongName": {
	// "code": "test-err:this-is-avery-long-name",
	// "description": "",
	// "protocol_code": 0
	// }
	// }
	// ]
}

func ExampleGetCaller() {
	_, caller := GetCaller(Internal, 1)
	fmt.Println(caller)
	// Output:
	// ExampleGetCaller
}

func ExampleNested() {
	data := &TestErrGroup{
		TestErr1:            &Def{Msg: `msg1`},
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`, Msg: `msg2`},
	}
	InitErrGroup(data)
	ze := data.TestErr1.Wrap(data.Err.Wrap(errors.New(`original-error`)))
	fmt.Println(ze.Error(), ze.CallerName, ze.Def.Code)

	ze = data.TestErr1.Wrapf(data.Err.Wrap(errors.New(`original-error`)), `wrap message: %s`, `ad`)
	fmt.Println(ze.Error(), ze.CallerName, ze.Def.Code)
	// Output:
	// msg1 <- msg2 <- original-error ExampleNested/ExampleNested test-err:test-err1
	// msg1 <- msg2 <- original-error ExampleNested/ExampleNested test-err:test-err1
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
		RespondMessage(true),
	)
	m.RegisterGroups()
	defer unregister()
	rsp := Internal.GetResponser(errors.New(`original error`))
	mared, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(mared))

	// Output:
	// {"A":"zerror:internal","Msg":"original error"}
}

func ExampleDefaultDef() {

	unregister()
	m := New(DefaultPCode(500))
	errDef := DefaultDef(`msg`)
	m.RegisterGroups()
	defer unregister()
	fmt.Printf("%+v\n", errDef)
	// Output:
	// &{Code: Msg:msg Description: PCode:500 extensions:map[]}
}

type ForDefaultGroup struct {
	Err1 *Def
}

func ExampleDefaultDef2() {

	unregister()
	d := &ForDefaultGroup{Err1: DefaultDef(`default`)}
	fmt.Println(d.Err1.PCode)
	m := New(DefaultPCode(500))
	m.RegisterGroups(d)
	defer unregister()

	fmt.Println(d.Err1.PCode)
	// Output:
	// -1
	// 500
}

func ExampleDef_Is() {

	originalError := errors.New(`original error`)
	def := DefaultDef(`default`)
	def1 := DefaultDef(`default2`)
	wrapped := def.Wrap(originalError)
	wrapped1 := def1.Wrap(wrapped)
	fmt.Println(def.Make(wrapped))
	fmt.Println(def1.Make(wrapped1))
	fmt.Println(def.Make(wrapped1))
	fmt.Println(def1.Make(wrapped))
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
