package gin_ze

import (
	"github.com/EchoUtopia/zerror"
	logrus_ze "github.com/EchoUtopia/zerror/logrus"
	"github.com/gin-gonic/gin"
)

const (
	LogWhenRespond = `log_when_respond`
)

func JSON(c *gin.Context, err error) {
	if !zerror.Registered() {
		panic(`groups not registered`)
	}
	var def *zerror.Def
	zerr, ok := err.(*zerror.Error)
	if !ok {
		def = zerror.InternalError
		location, caller := zerror.GetCaller(def, 2)
		zerr = zerror.NewError(err, def, ``, &zerror.ZContext{CallLocation: location, CallerName: caller})
	} else {
		def = zerr.Def
	}

	c.JSON(int(def.PCode), def.GetResponser(zerr))
	c.Abort()
	if _, logWhenRespond := zerror.Manager.GetExtension(LogWhenRespond); logWhenRespond {
		logrus_ze.LogCtx(c.Request.Context(), zerr)
	}
}
