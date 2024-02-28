package service

import (
	"sync"

	"github.com/kiga-hub/arc/logging"
)

// Service -
type Service struct {
	logger  logging.ILogger
	service *sync.Map
}

// New -
func New(opts ...Option) (*Service, error) {
	srv := loadOptions(opts...)
	srv.service = new(sync.Map) // Others
	return srv, nil
}

// Start -
func (s *Service) Start() {
	s.logger.Info("service started")
}

// Stop -
func (s *Service) Stop() {
}
