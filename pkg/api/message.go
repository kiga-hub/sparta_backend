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
		SetSummary("create websocket connectio").
		SetDescription(`create websocket conenction`).
		AddParamQuery("", "device", "app、web", false).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/health", s.Health).
		SetOperationId(`get app healthy`).
		SetSummary("get app healthy").
		SetDescription(`get app healthy`).
		AddResponse(http.StatusOK, ``, nil, nil)
}

// Connect -
func (s *Server) Connect(c echo.Context) error {

	s.logger.Info("websocket connect")

	return s.ws.HandleConnections(c)
}

// Health -
func (s *Server) Health(c echo.Context) error {

	resp := map[string]int{
		"health": 100,
	}
	return c.JSON(http.StatusOK, resp)
}
