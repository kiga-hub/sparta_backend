package api

import (
	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/websocket/pkg/service"
	"github.com/kiga-hub/websocket/pkg/ws"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Handler - 对外接口
type Handler interface {
	Setup(echoswagger.ApiRoot, string)
}

// Server - api处理器
type Server struct {
	logger logging.ILogger
	srv    *service.Service
	ws     *ws.WebsocketServer
}

// New - 初始化
func New(opts ...Option) Handler {
	srv := loadOptions(opts...)
	return srv
}
