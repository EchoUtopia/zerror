package zerror

const (
	CodeInternal = `zerror:internal`
)

var Internal = &Def{
	Code:        CodeInternal,
	Status:      StatusInternal,
	Msg:         `internal error`,
	Description: `server internal error`,
}

var BadRequest = &Def{
	Code:        `zerror:bad_request`,
	Status:      StatusBadRequest,
	Msg:         `bad request`,
	Description: `bad request`,
}

var Forbidden = &Def{
	Code:        `zerror:forbidden`,
	Status:      StatusPermissionDenied,
	Msg:         `forbidden`,
	Description: `you are forbidden to access`,
}

var NotFound = &Def{
	Code:        `zerror:not_found`,
	Msg:         `not found`,
	Status:      StatusNotFound,
	Description: `resource not found`,
}

var Unauthenticated = &Def{
	Code:        `zerror:unauthenticated`,
	Msg:         `unauthenticated`,
	Status:      StatusUnauthenticated,
	Description: `please login`,
}

var AlreadyExists = &Def{
	Code:        `zerror:already_exists`,
	Msg:         `already exists`,
	Status:      StatusAlreadyExists,
	Description: `already exists`,
}
