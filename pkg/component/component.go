package component

import (
	"context"

	"github.com/kiga-hub/websocket/pkg/service"
	"github.com/kiga-hub/websocket/pkg/ws"

	"github.com/davecgh/go-spew/spew"
	platformConf "github.com/kiga-hub/arc/conf"
	"github.com/kiga-hub/arc/logging"
	logConf "github.com/kiga-hub/arc/logging/conf"
	"github.com/kiga-hub/arc/micro"
	"github.com/kiga-hub/arc/micro/conf"
	"github.com/kiga-hub/websocket/pkg/api"
	"github.com/pangpanglabs/echoswagger/v2"
)

// WebScoketComponentElementKey is Element Key for WebScoketComponent
var WebScoketComponentElementKey = micro.ElementKey("WebScoketComponent")

// WebScoketComponent is Component for WebScoketComponent
type WebScoketComponent struct {
	micro.EmptyComponent
	config *conf.BasicConfig
	logger logging.ILogger
	api    api.Handler
	srv    *service.Service
	ws     *ws.WebsocketServer
}

// Name of the component
func (c *WebScoketComponent) Name() string {
	return "WebScoketComponent"
}

// PreInit called before Init()
func (c *WebScoketComponent) PreInit(ctx context.Context) error {
	service.SetDefaultConfig()
	return nil
}

// Init the component
func (c *WebScoketComponent) Init(server *micro.Server) (err error) {
	c.config = conf.GetBasicConfig()
	spew.Dump(c.config) // 打印基础配置信息

	// 获取日志接口
	elLogger := server.GetElement(&micro.LoggingElementKey)
	if elLogger != nil {
		c.logger = elLogger.(logging.ILogger)
		spew.Dump(logConf.GetLogConfig()) // 打印日志配置信息
	}

	// 初始化基础服务
	if c.srv, err = service.New(
		service.WithLogger(c.logger),
	); err != nil {
		return err
	}

	c.ws = ws.NewServer()
	// 初始化web api接口服务
	c.api = api.New(
		api.WithLogger(c.logger),
		api.WithService(c.srv),
		api.WithWebsocketServer(c.ws),
	)

	return nil
}

// SetDynamicConfig 加载动态配置回调
func (c *WebScoketComponent) SetDynamicConfig(nf *platformConf.NodeConfig) error {
	return nil
}

// OnConfigChanged 动态配置修改回调函数
func (c *WebScoketComponent) OnConfigChanged(*platformConf.NodeConfig) error {
	return micro.ErrNeedRestart
}

// SetupHandler 安装路由
func (c *WebScoketComponent) SetupHandler(root echoswagger.ApiRoot, base string) error {
	root.Echo().Static("/static", "./static")
	c.api.Setup(root, base)
	return nil
}

// Start the component
func (c *WebScoketComponent) Start(ctx context.Context) error {
	if c.srv != nil {
		go c.srv.Start()
	}
	return nil
}

// Stop the component
func (c *WebScoketComponent) Stop(ctx context.Context) error {
	if c.srv != nil {
		c.srv.Stop()
		// 移除所有websocket连接
		c.ws.RemoveAllConn()
	}
	return nil
}
