package zerror

const (
	BizCodeInternal        = `zerror:internal`
	BizCodeBadRequest      = `zerror:bad_request`
	BizCodeForbidden       = `zerror:forbidden`
	BizCodeNotFound        = `zerror:not_found`
	BizCodeUnauthenticated = `zerror:unauthenticated`
	BizCodeAlreadyExists   = `zerror:already_exists`
)

var InternalDef = &Def{
	Code:        BizCodeInternal,
	PCode:       CodeInternal,
	Msg:         `internal error`,
	Description: `this is server internal error, please contact admin`,
}

var BadRequestDef = &Def{
	Code:        BizCodeBadRequest,
	PCode:       CodeBadRequest,
	Msg:         `bad request`,
	Description: `check your request parameters`,
}

var ForbiddenDef = &Def{
	Code:        BizCodeForbidden,
	PCode:       CodePermissionDenied,
	Msg:         `forbidden`,
	Description: `you are forbidden to access`,
}

var NotFoundDef = &Def{
	Code:        BizCodeNotFound,
	Msg:         `not found`,
	PCode:       CodeNotFound,
	Description: `resource not found`,
}

var UnauthenticatedDef = &Def{
	Code:        BizCodeUnauthenticated,
	Msg:         `unauthenticated`,
	PCode:       CodeUnauthenticated,
	Description: `please login`,
}

var AlreadyExistsDef = &Def{
	Code:        BizCodeAlreadyExists,
	Msg:         `already exists`,
	PCode:       CodeAlreadyExists,
	Description: `already exists`,
}
