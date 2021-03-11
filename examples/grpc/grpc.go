package grpc

import (
	"context"
	"github.com/EchoUtopia/zerror"
	logrus_ze "github.com/EchoUtopia/zerror/examples/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var codeConvertMap = map[zerror.ProtocolCode]codes.Code{
	zerror.CodeOk: codes.OK,

	zerror.CodeBadRequest:         codes.InvalidArgument,
	zerror.CodeUnauthenticated:    codes.Unauthenticated,
	zerror.CodePermissionDenied:   codes.PermissionDenied,
	zerror.CodeNotFound:           codes.NotFound,
	zerror.CodeDeadlineExceeded:   codes.DeadlineExceeded,
	zerror.CodeAlreadyExists:      codes.AlreadyExists,
	zerror.CodeFailedPrecondition: codes.FailedPrecondition,
	zerror.CodeResourceExhausted:  codes.ResourceExhausted,
	zerror.CodeCancelled:          codes.Canceled,

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

//  you can set the business code and grpc code in the interceptor
type interceptor struct {
	logger logrus.FieldLogger
}

// your response message should have a Code field, which represents your business code
func (i *interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	rsp, err := handler(ctx, req)
	if err != nil {
		zerr, ok := err.(*zerror.Error)
		if ok {
			// if rsp has Code field
			err = status.Error(ConvertToGrpcCode(zerr.Def.PCode), zerr.Code)
		}
		logrus_ze.Log(err)
	}
	return rsp, err
}

type fakeClient struct{}

func (fc *fakeClient) Call(ctx context.Context) error {
	return nil
}

func ExampleHandleGrpcResponse() error {
	client := &fakeClient{}
	if err := client.Call(context.TODO()); err != nil {
		sts, ok := status.FromError(err)
		if ok {
			ok, zErr := zerror.FromCode(sts.Message())
			if ok {
				return zErr
			}
		}
		return err
	}
	return nil
}
