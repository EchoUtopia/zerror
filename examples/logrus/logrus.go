package logrus

import (
	"context"
	"errors"
	"github.com/EchoUtopia/zerror"
	"github.com/sirupsen/logrus"
	"log"
)

const (
	ExtLogLvl             = `log_level`
	ExtLogger             = `logger`
	ExtExtractDataFromCtx = `extract_from_ctx`
)

type ExtractDataFromCtx func(context.Context) zerror.Data

func LogCtx(ctx context.Context, err error) {
	data := zerror.Data{}

	logLevel := logrus.InfoLevel
	zerr := &zerror.Error{}
	l, n := ``, ``
	if ok := errors.As(err, &zerr); ok {
		l, n = zerr.GetCaller()
		lli, ok := zerr.Def.GetExtension(ExtLogLvl)
		if ok {
			logLevel = lli.(logrus.Level)
		}
		data = zerr.Data
	} else {
		l, n = zerror.GetCaller(nil, 3)
	}
	extractor, ok := zerror.Manager.GetExtension(ExtExtractDataFromCtx)
	if ok {
		f := extractor.(ExtractDataFromCtx)
		for k, v := range f(ctx) {
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
	zerr := &zerror.Error{}
	if ok := errors.As(err, &zerr); ok {
		l, n = zerr.GetCaller()
		data = zerr.Data
		lli, ok := zerr.Def.GetExtension(ExtLogLvl)
		if ok {
			logLevel = lli.(logrus.Level)
		}
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

	iLogger, ok := zerror.Manager.GetExtension(ExtLogger)
	logger, isLogger := iLogger.(logrus.FieldLogger)
	if zerror.Manager.DebugMode() {
		if !ok || !isLogger {
			log.Panicf(`manager extension: %s not exist or is not logrus.FieldLogger`, ExtLogger)
		}
	} else if !ok || !isLogger {
		log.Printf(`manager extension: %s not exist or is not logrus.FieldLogger`, ExtLogger)
		log.Printf(`data: %v, err: %s`, data, err)
		return
	}
	logger.WithFields(logrus.Fields(data)).WithError(err).Log(level)
}
