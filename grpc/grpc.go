package grpc

import (
	"github.com/EchoUtopia/zerror"
	"google.golang.org/grpc/codes"
)

var codeConvertMap = map[zerror.ProtocolCode]codes.Code{
	zerror.CodeOk: codes.OK,

	zerror.CodeBadRequest:         codes.InvalidArgument,
	zerror.CodeUnauthenticated:    codes.Unauthenticated,
	zerror.CodePermissionDenied:   codes.PermissionDenied,
	zerror.CodeNotFound:           codes.NotFound,
	zerror.CodeDeadlineExceeded:   codes.DeadlineExceeded,
	zerror.CodeFailedPrecondition: codes.FailedPrecondition,

	zerror.CodeInternal:      codes.Internal,
	zerror.CodeUnimplemented: codes.Unimplemented,
	zerror.CodeUnavailable:   codes.Unavailable,
}

func SetCustomRelation(zCode zerror.ProtocolCode, grpcCode codes.Code) {
	codeConvertMap[zCode] = grpcCode
}

func ConvertToGrpcCode(c zerror.ProtocolCode) codes.Code {
	return codeConvertMap[c]
}
