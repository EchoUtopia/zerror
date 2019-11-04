package zerror

import (
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
	RegisterGroups(data)
	fmt.Println(JsonDumpGroups(``))
	errGroups = nil
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
	fmt.Println(getCaller(logrus.InfoLevel, 1))
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
	fmt.Println(ze.Error(), ze.internal, ze.callerName, ze.def.Code)
	// Output:
	// msg1: msg2: original-error false ExampleNested/ExampleNested test-err:test-err1
}
