package main

import (
	"errors"
	"github.com/EchoUtopia/zerror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		Token: &zerror.Def{Code: ``, HttpCode: 401, Msg: `token invalid`, LogLevel: logrus.InfoLevel, Description: ``},

		// if zerr code is not empty, then it's the final code
		Expired: &zerror.Def{Code: `custom-expired`, HttpCode: 401, Msg: `token expired`, LogLevel: logrus.DebugLevel, Description: ``},
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
		err = Auth.Token.WrapAsInner(originalErr)
		zerror.JSON(c, err)
	case `straight`:
		SmsCode.JSON(c, originalErr)
	default:
		zerror.JSON(c, err)
	}
}

func main() {

	logger := logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	manager := zerror.New(zerror.Logger(logger))

	// error group must be registered
	manager.RegisterGroups(Common, Auth)

	r := gin.Default()
	r.GET(`/error`, ErrHandler)
	r.Run(`:8989`)

	// when access /error, response is : `{"code":"args","data":null,"msg":null}`, http code is 400
	// log is :
	// DEBU[0010] args err                                      call_location=/Users/echo/go/src/github.com/EchoUtopia/zerror/examples/main.go/56 caller=ErrHandler error="original error"

	// when access /error?type=original, response is : `{"code":"unkown","data":null,"msg":null}`, http code is 500
	// log is :
	// ERRO[0001] unkown error                                  caller=ErrHandler error="original error"

	// when access /error?type=internal, response is `{"code":"auth:token","data":null,"msg":null}`, http code is 500
	// log is :
	// INFO[0033] token invalid                                 caller=ErrHandler error="original error"

	// when access /error?type=straight, response is `{"code":"sms:code","data":null,"msg":null}`, http code is 500
	// log is :
	// ERRO[0064] sms code                                      caller=ErrHandler error="original error"
}
