package api

import (
	"github.com/pangpanglabs/echoswagger/v2"
)

// Setup register api
func (s *Server) Setup(root echoswagger.ApiRoot, base string) {
	if s.srv != nil {
		s.setupRESTfulApi(root, base)
		// s.setupWebsocket(root, base)
		s.setupStaticFiles(root, base)
	}
}
