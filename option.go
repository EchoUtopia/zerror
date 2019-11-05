package zerror

import "github.com/sirupsen/logrus"

type Options struct {
	wordConnector  string
	codeConnector  string
	RespondMessage bool
	logger         logrus.FieldLogger
	responseFunc   func() Responser
}

type Option func(*Options)

func WordConnector(wc string) Option {
	return func(options *Options) {
		options.wordConnector = wc
	}
}

func CodeConnector(cc string) Option {
	return func(options *Options) {
		options.codeConnector = cc
	}
}

func RespondMessage(respondMessage bool) Option {
	return func(options *Options) {
		options.RespondMessage = respondMessage
	}
}

func Logger(logger logrus.FieldLogger) Option {
	return func(options *Options) {
		options.logger = logger
	}
}

func WithResponser(rf func() Responser) Option {
	return func(options *Options) {
		options.responseFunc = rf
	}
}
