package zerror

type Options struct {
	wordConnector  string
	codeConnector  string
	respondMessage bool
	respondMsgSet  bool
	render         func() Render
	defaultStatus  Status
	debugMode      bool
	extensions     map[string]interface{}
}

type Option func(*Options)

func WithRender(rf func() Render, renderMsg bool) Option {
	return func(options *Options) {
		options.render = rf
		if renderMsg {
			options.respondMessage = renderMsg
			options.respondMsgSet = true
		}
	}
}

func DefaultStatus(status Status) Option {
	return func(options *Options) {
		options.defaultStatus = status
	}
}

func SetDebugMode(debug bool) Option {
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

func (m *Zmanager) GetExtension(key string) (interface{}, bool) {
	value, ok := m.extensions[key]
	return value, ok
}

func (o *Options) DebugMode() bool {
	return o.debugMode
}
