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

// //  you can set the business code and grpc code in the interceptor
// type interceptor struct {
// 	logger logrus.FieldLogger
// }
// // your response message should have a Code field, which represents your business code
// func (i *interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//
// 	h, err := handler(ctx, req)
// 	if err != nil {
// 		zerr, ok := err.(*zerror.Error)
// 		if ok {
// 			h.Code = zerr.Def.Code
// 			err = status.Error(grpc_ze.ConvertToGrpcCode(zerr.Def.PCode), zerr.Error())
// 		}
//      logrus_ze.Log(err)
// 	}
// 	return h, err
// }
