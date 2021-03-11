package zerror

type ProtocolCode int

const (
	CodeInvalid ProtocolCode = 0
	CodeOk      ProtocolCode = 200

	CodeBadRequest         ProtocolCode = 400
	CodeUnauthenticated    ProtocolCode = 401
	CodePermissionDenied   ProtocolCode = 403
	CodeNotFound           ProtocolCode = 404
	CodeDeadlineExceeded   ProtocolCode = 408
	CodeAlreadyExists      ProtocolCode = 409
	CodeFailedPrecondition ProtocolCode = 412
	CodeResourceExhausted  ProtocolCode = 429
	CodeCancelled          ProtocolCode = 499

	CodeInternal      ProtocolCode = 500
	CodeUnimplemented ProtocolCode = 501
	CodeUnavailable   ProtocolCode = 503
)
