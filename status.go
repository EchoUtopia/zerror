package zerror

type Status int

const (
	StatusInvalid Status = 0
	StatusOk      Status = 200

	StatusBadRequest         Status = 400
	StatusUnauthenticated    Status = 401
	StatusPermissionDenied   Status = 403
	StatusNotFound           Status = 404
	StatusDeadlineExceeded   Status = 408
	StatusAlreadyExists      Status = 409
	StatusFailedPrecondition Status = 412
	StatusResourceExhausted  Status = 429
	StatusCancelled          Status = 499

	StatusInternal      Status = 500
	StatusUnimplemented Status = 501
	StatusUnavailable   Status = 503
)
