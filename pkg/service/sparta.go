package service

import (
	"github.com/kiga-hub/websocket/pkg/models"
)

// CreatingParticles -
func (s *Service) CreatingParticles(sparta *models.Sparta) interface{} {

	result := sparta.ProcessSparta()
	s.logger.Info("Sparta: ", result)

	return result
}
