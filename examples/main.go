package main

import (
	"errors"
	"github.com/EchoUtopia/zerror"
	gin_ze "github.com/EchoUtopia/zerror/gin"
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

func HandleOriginal(c *gin.Context) {
	originalErr := errors.New(`original error`)
	gin_ze.JSON(c, originalErr)
}

func HandleInternal(c *gin.Context) {
	originalErr := errors.New(`original error`)
	err := zerror.InternalDef.Wrap(originalErr)
	err = Common.Args.Wrap(err)
	gin_ze.JSON(c, err)
}

func HandleDefault(c *gin.Context) {

	originalErr := errors.New(`original error`)
	err := Common.Args.Wrap(originalErr).WithData(zerror.Data{`custom key`: `custom value`})
	gin_ze.JSON(c, err)
}

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	manager := zerror.New(
		// zerror.DebugMode(true),
		// zerror.WithResponser(),
		// zerror.DefaultPCode(zerror.CodeBadRequest),
		zerror.RespondMessage(true),
		zerror.Extend(zerror.ExtLogger, logrus.StandardLogger()),
		zerror.Extend(gin_ze.LogWhenRespond, true),
	)

	// error group must be registered
	manager.RegisterGroups(Common, Auth)

	r := gin.Default()
	r.GET(`/error`, HandleDefault)
	r.GET(`/error/original`, HandleOriginal)
	r.GET(`/error/internal`, HandleInternal)
	r.Run(`:8989`)

	// just go to see the http response and the logs in server sever side
}
