package gants

// Option is the optional function.
type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// Options contains all options which will be applied when instantiating an gants pool.
type Options struct {
	EnableStatistic bool

	MaxWorkerCount int

	MaxTaskQueue int
}

func WithStatistic() Option {
	return func(opts *Options) {
		opts.EnableStatistic = true
	}
}

func WithMaxWorkerCount(maxWorkerCount int) Option {
	return func(opts *Options) {
		opts.MaxWorkerCount = maxWorkerCount
	}
}

func WithMaxTaskQueue(maxTaskQueue int) Option {
	return func(opts *Options) {
		opts.MaxTaskQueue = maxTaskQueue
	}
}
