package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// setupMsg setupMsg test
func (s *Server) setupMsg(root echoswagger.ApiRoot, base string) {
	g := root.Group("Websocket", base+"/ws")
	g.GET("", s.Connect).
		SetOperationId(`websocket connect`).
		SetSummary("创建websocket连接").
		SetDescription(`创建websocket连接`).
		AddParamQuery("", "device", "app、web", true).
		AddResponse(http.StatusOK, ``, nil, nil)
}

// Connect -
func (s *Server) Connect(c echo.Context) error {

	s.logger.Info("websocket connect")

	return s.ws.HandleConnections(c)
}
