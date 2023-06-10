package config

type configOptionParam struct {
	fileChangeMonitor bool
}

type Option func(*configOptionParam)

type Options []Option

func (s Options) Do(i *configOptionParam) {
	for _, option := range s {
		option(i)
	}
}

func WithFileChangeMonitor() Option {
	return func(c *configOptionParam) {
		c.fileChangeMonitor = true
	}
}
