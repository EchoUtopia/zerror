package logrus_ze

import (
	"context"
	"github.com/EchoUtopia/zerror"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	ExtractDataFromCtxFunc ExtractDataFromCtx
)

type ExtractDataFromCtx func(context.Context) zerror.Data

func LogCtx(ctx context.Context, err error) {
	data := zerror.Data{}

	logLevel := logrus.ErrorLevel
	zerr, ok := err.(*zerror.Error)
	l, n := ``, ``
	if ok {
		l, n = zerr.GetCaller()
		li, ok := zerr.Def.GetExtension(zerror.ExtLogLvl)
		if ok {
			logLevel = li.(logrus.Level)
		}
		data = zerr.Data
	} else {
		l, n = zerror.GetCaller(nil, 2)
	}
	if ExtractDataFromCtxFunc != nil {
		for k, v := range ExtractDataFromCtxFunc(ctx) {
			data[k] = v
		}
	}
	data[`caller`] = n
	if l != `` {
		data[`call_location`] = l
	}
	getAndLog(err, data, logLevel)
}

func Log(err error) {
	data := zerror.Data{}
	logLevel := logrus.ErrorLevel
	l, n := ``, ``
	zerr, ok := err.(*zerror.Error)
	if ok {
		l, n = zerr.GetCaller()
		data = zerr.Data
		li, ok := zerr.Def.GetExtension(zerror.ExtLogLvl)
		if ok {
			logLevel = li.(logrus.Level)
		}
		data = zerr.Data
	} else {
		l, n = zerror.GetCaller(nil, 2)
	}
	data[`caller`] = n
	if l != `` {
		data[`call_location`] = l
	}
	getAndLog(err, data, logLevel)
}

func getAndLog(err error, data zerror.Data, level logrus.Level) {

	iLogger, ok := zerror.Manager.GetExtension(zerror.ExtLogLvl)
	logger, isLogger := iLogger.(logrus.FieldLogger)
	if zerror.Manager.DebugMode() {
		if !ok || !isLogger {
			log.Panicf(`manager extension: %s not exist or is not logrus.FieldLogger`, zerror.ExtLogLvl)
		}
	} else if !ok || !isLogger {
		log.Printf(`manager extension: %s not exist or is not logrus.FieldLogger`, zerror.ExtLogLvl)
		log.Printf(`data: %v, err: %s`, data, err)
		return
	}
	logger.WithFields(logrus.Fields(data)).WithError(err).Log(level)
}
