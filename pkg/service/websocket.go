package service

import (
	"github.com/gorilla/websocket"
	"github.com/kiga-hub/sparta_backend/pkg/ws"
	"github.com/labstack/echo/v4"
)

// WSCompute - WSCompute
func (s *Service) WSCompute(conn *websocket.Conn, c echo.Context) error {
	client := ws.NewConn(conn, s.logger)

	go client.ComputeReadLoop()
	go client.ComputeWriteLoop()

	return nil
}

// WSConvertToVisualFormat -
func (s *Service) WSConvertToVisualFormat(conn *websocket.Conn, c echo.Context) error {
	client := ws.NewConn(conn, s.logger)

	go client.ConvertToVisualFormatReadLoop()
	go client.ConvertToVisualFormatWriteLoop()

	return nil
}
