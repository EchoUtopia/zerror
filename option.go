package zerror

type Options struct {
	wordConnector  string
	codeConnector  string
	RespondMessage bool
	responseFunc   func() Responser
	defaultPCode   ProtocolCode
	debugMode      bool
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

func WithResponser(rf func() Responser) Option {
	return func(options *Options) {
		options.responseFunc = rf
	}
}

func DefaultPCode(code ProtocolCode) Option {
	return func(options *Options) {
		options.defaultPCode = code
	}
}

func DebugMode(debug bool) Option {
	return func(options *Options) {
		options.debugMode = debug
	}
}
