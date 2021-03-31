module github.com/EchoUtopia/zerror/examples/v2

go 1.13

require (
	github.com/EchoUtopia/zerror/v2 v2.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.31.0
)

replace github.com/EchoUtopia/zerror/v2 => ../
