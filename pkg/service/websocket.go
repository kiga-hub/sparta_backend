package service

import (
	"github.com/gorilla/websocket"
	"github.com/kiga-hub/sparta_backend/pkg/ws"
	"github.com/labstack/echo/v4"
)

// WSConnect - WSConnect
func (s *Service) WSConnect(conn *websocket.Conn, c echo.Context) error {
	client := ws.NewConn(conn, s.logger)

	go client.ReadLoop()
	go client.WriteLoop()

	return nil
}
