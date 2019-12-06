package logrus_ze

import (
	"context"
	"github.com/EchoUtopia/zerror"
	"github.com/sirupsen/logrus"
)

var (
	Logger                 logrus.FieldLogger
	ExtractDataFromCtxFunc ExtractDataFromCtx
)

type ExtractDataFromCtx func(context.Context) zerror.Data

func Log(ctx context.Context, def *zerror.Def, err error) {
	data := zerror.Data{}
	l, n := zerror.GetCaller(def, 2)
	data[`caller`] = n
	if l != `` {
		data[`call_location`] = l
	}
	logLevel := logrus.ErrorLevel
	zerr, ok := err.(*zerror.Error)
	if ok {
		for k, v := range zerr.Data {
			data[k] = v
		}
		li, ok := zerr.Def.Extensions[zerror.ExtLogLvl]
		if ok {
			logLevel = li.(logrus.Level)
		}
	}
	if ExtractDataFromCtxFunc != nil {
		for k, v := range ExtractDataFromCtxFunc(ctx) {
			data[k] = v
		}
	}
	Logger.WithFields(logrus.Fields(data)).WithError(err).Log(logLevel, def.Msg)
}
