package zerror

type Options struct {
	wordConnector  string
	codeConnector  string
	RespondMessage bool
	RespondMsgSet  bool
	responseFunc   func() Responser
	defaultPCode   ProtocolCode
	debugMode      bool
	extensions     map[string]interface{}
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
		options.RespondMsgSet = true
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

func Extend(key string, value interface{}) Option {
	return func(options *Options) {
		if options.extensions == nil {
			options.extensions = make(map[string]interface{})
		}
		options.extensions[key] = value
	}
}

func (o *Options) GetExtension(key string) (interface{}, bool) {
	value, ok := Manager.extensions[key]
	return value, ok
}

func (o *Options) DebugMode() bool {
	return o.debugMode
}
