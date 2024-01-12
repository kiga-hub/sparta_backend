package api

import (
	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/websocket/pkg/service"
)

// Option is a function that will set up option.
type Option func(opts *Server)

func loadOptions(options ...Option) *Server {
	opts := &Server{}
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
	return func(opts *Server) {
		opts.logger = logger
	}
}

// WithService -
func WithService(g *service.Service) Option {
	return func(opts *Server) {
		opts.srv = g
	}
}
