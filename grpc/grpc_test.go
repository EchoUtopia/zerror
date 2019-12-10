package grpc

import (
	"github.com/EchoUtopia/zerror"
	"testing"
)

func BenchmarkConvertToGrpcCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConvertToGrpcCode(zerror.CodeInternal)
	}
}