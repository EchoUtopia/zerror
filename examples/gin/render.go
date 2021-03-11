package gin

import (
	"errors"
	"github.com/EchoUtopia/zerror"
	logrus_ze "github.com/EchoUtopia/zerror/examples/logrus"
	"github.com/gin-gonic/gin"
)

const (
	ExtLogWhenRespond = `log_when_respond`
)

func JSON(c *gin.Context, err error) {
	if !zerror.Manager.Registered() {
		panic(`groups not registered`)
	}
	var zerr *zerror.Error
	if ok := errors.As(err, &zerr); !ok {
		zerr = zerror.Internal.Wrap(err)
	}
	c.JSON(int(zerr.PCode), zerr.Render())
	c.Abort()
	if _, logWhenRespond := zerror.Manager.GetExtension(ExtLogWhenRespond); logWhenRespond {
		logrus_ze.LogCtx(c.Request.Context(), err)
	}
}
