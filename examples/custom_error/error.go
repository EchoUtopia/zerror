package custom_error

import "github.com/EchoUtopia/zerror/v2"

var (
	Auth = &auth{
		Prefix: "auth",
		// the Args def will be initialized automatically
	}
)

type auth struct {
	Prefix  string
	Token   *zerror.Def
	Expired *zerror.Def
}
