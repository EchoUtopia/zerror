package main

import (
	"context"
	"errors"
	"github.com/EchoUtopia/zerror/examples/v2/custom_error"
	gin_ze "github.com/EchoUtopia/zerror/examples/v2/gin"
	logrus_ze "github.com/EchoUtopia/zerror/examples/v2/logrus"
	"github.com/EchoUtopia/zerror/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


// you can use middleware to extract some data in context, and log them
func SetCtxValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, `context key`, `context value`))
	}
}

func HandleOriginal(c *gin.Context) {
	originalErr := errors.New(`original error`)
	gin_ze.JSON(c, originalErr)
}

func HandleInternal(c *gin.Context) {
	originalErr := errors.New(`original error`)
	err1 := zerror.Internal.Wrap(originalErr)
	err := custom_error.Common.Args.Wrap(err1)
	gin_ze.JSON(c, err)
}

func HandleDefault(c *gin.Context) {

	originalErr := errors.New(`original error`)
	err := custom_error.Auth.Token.Wrap(originalErr).WithData(zerror.Data{`custom key`: `custom value`})
	gin_ze.JSON(c, err)
}

func ExtractFromCtx(ctx context.Context) zerror.Data {
	out := make(zerror.Data)
	value, ok := ctx.Value(`context key`).(string)
	if ok {
		out[`context key`] = value
	}
	return out
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	manager := zerror.Init(
		// zerror.DebugMode(true),
		zerror.DefaultStatus(zerror.StatusBadRequest),
		zerror.Extend(logrus_ze.ExtLogger, logrus.StandardLogger()),
		zerror.Extend(gin_ze.ExtLogWhenRespond, true),
		zerror.Extend(logrus_ze.ExtExtractDataFromCtx, logrus_ze.ExtractDataFromCtx(ExtractFromCtx)),
	)

	// error group must be registered
	manager.RegisterGroups(custom_error.Common, custom_error.Auth)

	r := gin.Default()
	r.Use(SetCtxValue())
	r.GET(`/error`, HandleDefault)
	r.GET(`/error/original`, HandleOriginal)
	r.GET(`/error/internal`, HandleInternal)
	r.GET(`/errors`, func(c *gin.Context) {
		c.JSON(200, gin.H{
			`code`: `ok`,
			`data`: manager.GetErrorGroups(),
		})
	})
	r.Run(`:8989`)

	// just go to see the http response and the logs in server sever side
}
