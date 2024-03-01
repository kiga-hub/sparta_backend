package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// setupWebsocket -
func (s *Server) setupWebsocket(root echoswagger.ApiRoot, base string) {
	g := root.Group("Websocket", base+"/ws")

	g.GET("", s.WsConnect).
		SetOperationId(`websocket connect`).
		SetSummary("create websocket conenction").
		SetDescription(`create websocket conenction`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/health", s.Health).
		SetOperationId(`get app healthy`).
		SetSummary("get app healthy").
		SetDescription(`get app healthy`).
		AddResponse(http.StatusOK, ``, nil, nil)

}

// WsConnect - WsConnect
func (s *Server) WsConnect(c echo.Context) error {

	var socketUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := socketUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrSocketConnectFailCode, utils.ErrSocketConnectFailMsg, err))
	}

	// print remote address
	fmt.Println("client connected:", conn.RemoteAddr())

	if err := s.srv.WSConnect(conn, c); err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrSocketConnectFailCode, utils.ErrSocketConnectFailMsg, err))
	}

	return nil
}

// Health -
func (s *Server) Health(c echo.Context) error {

	resp := map[string]int{
		"health": 100,
	}
	return c.JSON(http.StatusOK, resp)
}
