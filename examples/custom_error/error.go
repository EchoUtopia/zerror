package custom_error

import (
	logrus_ze "github.com/EchoUtopia/zerror/examples/v2/logrus"
	"github.com/EchoUtopia/zerror/v2"
	"github.com/sirupsen/logrus"
)

var (
	Auth = &auth{
		Prefix: "auth",
		// the Args def will be initialized automatically
	}

	Common = &CommonGroup{
		// if prefix is empty, there's no prefix in code
		Prefix: "",

		// the code will be `args`
		Args: (&zerror.Def{Code: ``, Status: 400, Msg: `args err`, Description: ``}).Extend(logrus_ze.ExtLogLvl, logrus.DebugLevel),
	}

	SmsCode           = &zerror.Def{Code: `sms:code`, Status: 500, Msg: `sms code`, Description: ``}
)

// this is error group
type CommonGroup struct {
	Prefix string
	Args   *zerror.Def
}

type auth struct {
	Prefix  string
	Token   *zerror.Def
	Expired *zerror.Def
}
