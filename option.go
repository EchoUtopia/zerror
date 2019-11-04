package zerror

import "github.com/sirupsen/logrus"

type Options struct {
	wordConnector string
	codeConnector string
	debug         bool
	logger        logrus.FieldLogger
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

func Debug(debug bool) Option {
	return func(options *Options) {
		options.debug = debug
	}
}

func Logger(logger logrus.FieldLogger) Option {
	return func(options *Options) {
		options.logger = logger
	}
}
