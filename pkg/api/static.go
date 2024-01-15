package api

import (
	"github.com/kiga-hub/websocket/pkg/service"
	"github.com/pangpanglabs/echoswagger/v2"
)

// setupStaticFiles -
func (s *Server) setupStaticFiles(root echoswagger.ApiRoot, base string) {
	g := root.Group("StaticFiles", base+"/static")
	// Access static resource files
	g.EchoGroup().Static("", service.GetConfig().DataDir)
}
