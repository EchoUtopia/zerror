module github.com/EchoUtopia/zerror/examples/v2

go 1.13

require (
	github.com/EchoUtopia/zerror/examples v0.0.0-20210313062224-47598bd95585
	github.com/EchoUtopia/zerror/v2 v2.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.31.0
)

replace github.com/EchoUtopia/zerror/v2 => ../
