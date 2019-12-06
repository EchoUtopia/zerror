package main

import (
	"errors"
	"github.com/EchoUtopia/zerror"
	gin_ze "github.com/EchoUtopia/zerror/gin"
	"github.com/EchoUtopia/zerror/logrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	Common = &CommonGroup{
		// if prefix is empty, there's no prefix in code
		Prefix: "",

		// the code will be `args`
		Args: (&zerror.Def{Code: ``, PCode: 400, Msg: `args err`, Description: ``}).Extend(zerror.ExtLogLvl, logrus.DebugLevel),
	}

	Auth = &AuthGroup{
		// prefix in code, seperated by ':"
		Prefix: "auth",

		// the code will be `auth:token`
		Token: zerror.DefaultDef(`invalid token`),

		// if zerr code is not empty, then it's the final code
		Expired: zerror.DefaultDef(`token expired`).WithCode(`custom-expired`),
	}

	SmsCode = &zerror.Def{Code: `sms:code`, PCode: 500, Msg: `sms code`, Description: ``}
)

// this is error group
type CommonGroup struct {
	Prefix string
	Args   *zerror.Def
}

type AuthGroup struct {
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
		gin_ze.JSON(c, originalErr)
	case `internal`:
		err = zerror.InternalError.Wrap(err)
		gin_ze.JSON(c, err)
	case `straight`:
		err = SmsCode.Wrap(originalErr)
		gin_ze.JSON(c, err)
	case `undefined`:
		gin_ze.JSON(c, originalErr)
	default:
		gin_ze.JSON(c, err)
	}
}

func main() {

	logrus_ze.Logger = logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	manager := zerror.New(
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
