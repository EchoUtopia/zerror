package grpc

import (
	"context"
	"errors"
	"github.com/EchoUtopia/zerror/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var zStatusToGCode = map[zerror.Status]codes.Code{
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

var gCodeToZStatus = map[codes.Code]zerror.Status{}

func init(){
	for k, v := range zStatusToGCode {
		gCodeToZStatus[v] = k
	}
}

func SetCustomRelation(status zerror.Status, grpcCode codes.Code) {
	zStatusToGCode[status] = grpcCode
	gCodeToZStatus[grpcCode] = status
}

// convert zerror status to grpc code
func ZStatus2GCode(s zerror.Status) codes.Code {
	code, ok := zStatusToGCode[s]
	if ok {
		return code
	}
	return codes.Unknown
}

// convert grpc code to zerror status
func GCode2ZStatus(c codes.Code)zerror.Status{
	return gCodeToZStatus[c]
}

//  you can set the business code and grpc code in the interceptor
type interceptor struct {}

// your response message should have a Code field, which represents your business code
func (i *interceptor) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	rsp, err := handler(ctx, req)
	if err != nil {
		var zerr *zerror.Error
		if errors.As(err, &zerr) {
			err = status.Error(ZStatus2GCode(zerr.Def.Status), zerr.Error())
		}
		// logrus_ze.Log(err)
	}
	return rsp, err
}
