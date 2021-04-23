package zerror

type defMapT map[string]*Def

var (
	defMap = defMapT{}
)

func init() {
	defMap.init()
}

const (
	CodeInternal        = `zerror:internal`
	codeBadRequest      = `zerror:bad_request`
	codeForbidden       = `zerror:forbidden`
	codeNofFound        = `zerror:not_found`
	codeUnauthenticated = `zerror:unauthenticated`
	codeAlreadyExists   = `zerror:already_exists`
)

var Internal = &Def{
	Code:        CodeInternal,
	Status:      StatusInternal,
	Msg:         `internal error`,
	Description: `server internal error`,
}

var BadRequest = &Def{
	Code:        codeBadRequest,
	Status:      StatusBadRequest,
	Msg:         `bad request`,
	Description: `bad request`,
}

var Forbidden = &Def{
	Code:        codeForbidden,
	Status:      StatusPermissionDenied,
	Msg:         `forbidden`,
	Description: `you are forbidden to access`,
}

var NotFound = &Def{
	Code:        codeNofFound,
	Msg:         `not found`,
	Status:      StatusNotFound,
	Description: `resource not found`,
}

var Unauthenticated = &Def{
	Code:        codeUnauthenticated,
	Msg:         `unauthenticated`,
	Status:      StatusUnauthenticated,
	Description: `please login`,
}

var AlreadyExists = &Def{
	Code:        codeAlreadyExists,
	Msg:         `already exists`,
	Status:      StatusAlreadyExists,
	Description: `already exists`,
}

func (m *defMapT) init() {
	inited := defMapT{
		CodeInternal:        Internal,
		codeBadRequest:      BadRequest,
		codeForbidden:       Forbidden,
		codeNofFound:        NotFound,
		codeUnauthenticated: Unauthenticated,
		codeAlreadyExists:   AlreadyExists,
	}
	*m = inited
}
