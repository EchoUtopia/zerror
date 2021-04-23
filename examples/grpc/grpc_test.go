package grpc

import (
	"github.com/EchoUtopia/zerror/v2"
	"testing"
)

func BenchmarkConvertToGrpcCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ZStatus2GCode(zerror.StatusInternal)
	}
}
