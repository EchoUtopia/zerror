package grpc

import (
	"context"
	logrus_ze "github.com/EchoUtopia/zerror/examples/v2/logrus"
	"github.com/EchoUtopia/zerror/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var codeConvertMap = map[zerror.Status]codes.Code{
	zerror.StatusOk: codes.OK,

	zerror.StatusBadRequest:         codes.InvalidArgument,
	zerror.StatusUnauthenticated:    codes.Unauthenticated,
	zerror.StatusPermissionDenied:   codes.PermissionDenied,
	zerror.StatusNotFound:           codes.NotFound,
	zerror.StatusDeadlineExceeded:   codes.DeadlineExceeded,
	zerror.StatusAlreadyExists:      codes.AlreadyExists,
	zerror.StatusFailedPrecondition: codes.FailedPrecondition,
	zerror.StatusResourceExhausted:  codes.ResourceExhausted,
	zerror.StatusCancelled:          codes.Canceled,

	zerror.StatusInternal:      codes.Internal,
	zerror.StatusUnimplemented: codes.Unimplemented,
	zerror.StatusUnavailable:   codes.Unavailable,
}

func SetCustomRelation(zCode zerror.Status, grpcCode codes.Code) {
	codeConvertMap[zCode] = grpcCode
}

func ConvertToGrpcCode(c zerror.Status) codes.Code {
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
			err = status.Error(ConvertToGrpcCode(zerr.Def.Status), zerr.Code)
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
			zErr, ok := zerror.FromCode(sts.Message())
			if ok {
				return zErr
			}
		}
		return err
	}
	return nil
}
