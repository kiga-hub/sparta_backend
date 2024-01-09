package service

import (
	"sync"

	"github.com/kiga-hub/arc/logging"
)

// Service - 服务结构
type Service struct {
	logger  logging.ILogger
	service *sync.Map
}

// New  - 初始化结构
func New(opts ...Option) (*Service, error) {
	srv := loadOptions(opts...)

	srv.service = new(sync.Map) // Others
	return srv, nil
}

// Start -
func (s *Service) Start() {
	s.logger.Info("service started")

}

// Stop - 停止服务
func (s *Service) Stop() {

}
