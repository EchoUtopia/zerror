package zerror

type ProtocolCode int

const (
	CodeInvalid ProtocolCode = -1
	CodeOk      ProtocolCode = 200

	CodeBadRequest         ProtocolCode = 400
	CodeUnauthenticated    ProtocolCode = 401
	CodePermissionDenied   ProtocolCode = 403
	CodeNotFound           ProtocolCode = 404
	CodeDeadlineExceeded   ProtocolCode = 408
	CodeFailedPrecondition ProtocolCode = 412

	CodeInternal      ProtocolCode = 500
	CodeUnimplemented ProtocolCode = 501
	CodeUnavailable   ProtocolCode = 503
)
