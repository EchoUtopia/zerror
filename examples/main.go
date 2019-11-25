package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/target-digital-transformation/kit/zerror"
)

var (
	Common = &CommonErr{
		// if prefix is empty, there's no prefix in code
		Prefix: "",

		// the code will be `args`
		Args: &zerror.Def{Code: ``, HttpCode: 400, Msg: `args err`, LogLevel: logrus.DebugLevel, Description: ``},
	}

	Auth = &AuthErr{
		// prefix in code, seperated by ':"
		Prefix: "auth",

		// the code will be `auth:token`
		Token: zerror.DefaultDef(`invalid token`),

		// if zerr code is not empty, then it's the final code
		Expired: zerror.DefaultDef(`token expired`).SetCode(`custom-expired`),
	}

	SmsCode = &zerror.Def{Code: `sms:code`, HttpCode: 500, Msg: `sms code`, LogLevel: logrus.ErrorLevel, Description: ``}
)

// this is error group
type CommonErr struct {
	Prefix string
	Args   *zerror.Def
}

type AuthErr struct {
	Prefix  string
	Token   *zerror.Def
	Expired *zerror.Def
}

func ErrHandler(c *gin.Context) {
	originalErr := errors.New(`original error`)

	err := Common.Args.Wrap(originalErr)

	errType := c.Query(`type`)
	switch errType {
	case `original`:
		zerror.JSON(c, originalErr)
	case `internal`:
		err = zerror.InternalError.Wrap(err)
		zerror.JSON(c, err)
	case `straight`:
		SmsCode.JSON(c, originalErr)
	case `undefined`:
		zerror.JSON(c, originalErr)
	default:
		zerror.JSON(c, err)
	}
}

func main() {

	logger := logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	manager := zerror.New(
		zerror.Logger(logger),
		zerror.RespondMessage(false),
	)

	// error group must be registered
	manager.RegisterGroups(Common, Auth)

	r := gin.Default()
	r.GET(`/error`, ErrHandler)
	r.Run(`:8989`)

	// when access /error, response is : `{"code":"args","data":null}`, http code is 400
	// log is :
	// DEBU[0010] args err                                      call_location=/Users/echo/go/src/github.com/EchoUtopia/zerror/examples/main.go/56 caller=ErrHandler error="original error"

	// when access /error?type=original, response is : `{"code":"unkown","data":null}`, http code is 500
	// log is :
	// ERRO[0001] unkown error                                  caller=ErrHandler error="original error"

	// when access /error?type=internal, response is `{"code":"zerror:internal","data":null}`, http code is 500
	// log is :
	// ERRO[0002] internal error                                call_location=/Users/echo/go/src/github.com/EchoUtopia/zerror/examples/main.go/48 caller=ErrHandler/ErrHandler error="args err: original error"

	// when access /error?type=straight, response is `{"code":"sms:code","data":null}`, http code is 500
	// log is :
	// ERRO[0064] sms code                                      caller=ErrHandler error="original error"

	// when access /error?type=undefined, response is `{"code":"zerror:undefined","data":null}`, http code is 500
	// log is :
	// ERRO[0006] unkown error                                  caller=ErrHandler error="original error"
}
