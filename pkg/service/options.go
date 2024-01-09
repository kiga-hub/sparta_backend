package service

import (
	"github.com/kiga-hub/arc/logging"
)

// Option is a function that will set up option.
type Option func(opts *Service)

func loadOptions(options ...Option) *Service {
	opts := &Service{}
	for _, option := range options {
		option(opts)
	}
	if opts.logger == nil {
		opts.logger = new(logging.NoopLogger)
	}
	return opts
}

// WithLogger -
func WithLogger(logger logging.ILogger) Option {
	return func(opts *Service) {
		opts.logger = logger
	}
}
