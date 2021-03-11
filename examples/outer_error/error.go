package outer_error

import "github.com/EchoUtopia/zerror"

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
