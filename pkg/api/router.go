package api

import (
	"github.com/pangpanglabs/echoswagger/v2"
)

// Setup register api
func (s *Server) Setup(root echoswagger.ApiRoot, base string) {
	if s.srv != nil {
		s.setupMsg(root, base)
	}
}
