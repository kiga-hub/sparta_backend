package api

import (
	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/sparta_backend/pkg/service"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Handler - for api
type Handler interface {
	Setup(echoswagger.ApiRoot, string)
}

// Server - api server
type Server struct {
	logger logging.ILogger
	srv    *service.Service
}

// New - create api server
func New(opts ...Option) Handler {
	srv := loadOptions(opts...)
	return srv
}
