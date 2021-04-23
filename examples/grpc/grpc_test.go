package grpc

import (
	"github.com/EchoUtopia/zerror/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"testing"
)

func BenchmarkConvertToGrpcCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ZStatus2GCode(zerror.StatusInternal)
	}
}

func TestZStatus2GCode(t *testing.T){
	code := ZStatus2GCode(zerror.StatusBadRequest)
	require.Equal(t, code, codes.InvalidArgument)
	code = ZStatus2GCode(zerror.Status(1000))
	require.Equal(t, code, codes.Unknown)
}

func TestGCode2ZStatus(t *testing.T) {
	status := GCode2ZStatus(codes.InvalidArgument)
	require.Equal(t, status, zerror.StatusBadRequest)
	status = GCode2ZStatus(codes.Unknown)
	require.Equal(t, status, zerror.StatusInvalid)
}

