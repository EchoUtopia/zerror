package zerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
	Manager = Init()
	m.Run()
}

func TestGenerateCode(t *testing.T) {
	data := &TestErr{
		TestErr1:            new(Def),
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`},
	}
	initErrGroup(data)
	require.Equal(t, `test-err:test-err1`, data.TestErr1.Code)
	require.Equal(t, `test-err:this-is-avery-long-name`, data.ThisISAVeryLongName.Code)

	require.Equal(t, `custom-code`, data.Err.Code)
	data.Err.Code = ``
	defMap.init()
	initErrGroup(data)
	require.Equal(t, `test-err:err`, data.Err.Code)

	data1 := &TestErr1{
		Err:    new(Def),
		Prefix: "",
	}
	initErrGroup(data1)
	require.Equal(t, `err`, data1.Err.Code)

	data1.Prefix = `custom-prefix`
	data1.Err.Code = ``
	initErrGroup(data1)
	require.Equal(t, `custom-prefix:err`, data1.Err.Code)
	require.Equal(t, Status(500), data1.Err.Status)

}

func ExampleGetCaller() {
	_, caller := GetCaller(Internal, 1)
	fmt.Println(caller)
	// Output:
	// ExampleGetCaller
}

func ExampleNested() {
	unregister()
	data := &TestErr{
		TestErr1:            &Def{Msg: `msg1`},
		ThisISAVeryLongName: new(Def),
		Err:                 &Def{Code: `custom-code`, Msg: `msg2`},
	}
	initErrGroup(data)
	ze := data.TestErr1.Wrap(data.Err.Wrap(errors.New(`original-error`)))
	fmt.Println(ze.Error())
	fmt.Println(ze.callerName, ze.Def.Code)

	ze = data.TestErr1.Wrapf(data.Err.Wrap(errors.New(`original-error`)), `wrap message: %s`, `ad`)
	fmt.Println(ze.Error())
	fmt.Println(ze.callerName, ze.Def.Code)
	// Output:
	// test-err:test-err1 | custom-code | original-error
	//ExampleNested/ExampleNested test-err:test-err1
	//test-err:test-err1(wrap message: ad) | custom-code | original-error
	//ExampleNested/ExampleNested test-err:test-err1
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
func (c *customeRsp) Error() string {
	return c.Msg
}

func ExampleCustomResponser() {
	unregister()
	m := Init(
		WithRender(func() Render {
			return new(customeRsp)
		}, true),
	)
	m.RegisterGroups()
	defer unregister()
	rsp := Internal.WithMsg(`original msg`).Render()
	mared, err := json.Marshal(rsp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(mared))

	// Output:
	// {"A":"zerror:internal","Msg":"zerror:internal(original msg)"}
}

func ExampleDefaultDef() {

	unregister()
	m := Init(DefaultStatus(500))
	data := &TestErr{}
	m.RegisterGroups(data)
	defer unregister()
	fmt.Printf("%+v\n", data.Err)
	// Output:
	// &{Code:test-err:err Msg: Description: Status:500 extensions:map[]}
}

func ExampleDef_Is() {

	originalError := errors.New(`original error`)
	def := &Def{Code: `def`}
	def1 := &Def{Code: `def1`}

	wrapped := def.Wrap(originalError)
	wrapped1 := def1.Wrap(wrapped)

	fmt.Println(wrapped1.Error())

	fmt.Println(def.Cause(wrapped))
	fmt.Println(def1.Cause(wrapped1))

	fmt.Println(def.Cause(wrapped1))
	fmt.Println(def1.Cause(wrapped))
	// Output:
	// def1 | def | original error
	//true
	//true
	//true
	//false
}

func BenchmarkDef_Wrap(b *testing.B) {

	originalError := errors.New(`original error`)
	def := &Def{Msg: `default`}
	for i := 0; i < b.N; i++ {
		def.wrapf(originalError, 1, ``)
	}
}

func TestData(t *testing.T) {
	zerr := Internal.New().
		WithData(map[string]interface{}{
			`a1`: `a2`,
		}).
		WithData(map[string]interface{}{
			`b1`: `b2`,
		}).
		WithKVs(`c1`, `c2`, `d1`, `d2`)
	expected := Data{
		`a1`: `a2`,
		`b1`: `b2`,
		`c1`: `c2`,
		`d1`: `d2`,
	}
	if !reflect.DeepEqual(zerr.Data, expected) {
		t.Fatal()
	}
}

func TestFromCode(t *testing.T) {
	zerr, ok := FromCode(BadRequest.Code)
	require.Equal(t, true, ok)
	require.Equal(t, zerr.Code, BadRequest.Code)
}
