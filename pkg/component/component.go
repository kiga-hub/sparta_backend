package component

import (
	"context"

	"github.com/kiga-hub/sparta_backend/pkg/service"
	"github.com/kiga-hub/sparta_backend/pkg/ws"

	"github.com/davecgh/go-spew/spew"
	platformConf "github.com/kiga-hub/arc/conf"
	"github.com/kiga-hub/arc/logging"
	logConf "github.com/kiga-hub/arc/logging/conf"
	"github.com/kiga-hub/arc/micro"
	"github.com/kiga-hub/arc/micro/conf"
	"github.com/kiga-hub/sparta_backend/pkg/api"
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
	// ws     *ws.WebsocketServer
}

// Name of the component
func (c *WebScoketComponent) Name() string {
	return "WebScoketComponent"
}

// PreInit called before Init()
func (c *WebScoketComponent) PreInit(ctx context.Context) error {
	service.SetDefaultConfig()
	ws.SetDefaultConfig()
	return nil
}

// Init the component
func (c *WebScoketComponent) Init(server *micro.Server) (err error) {
	c.config = conf.GetBasicConfig()
	spew.Dump(c.config) // print config

	// get logger
	elLogger := server.GetElement(&micro.LoggingElementKey)
	if elLogger != nil {
		c.logger = elLogger.(logging.ILogger)
		spew.Dump(logConf.GetLogConfig())
	}

	// init basic service
	if c.srv, err = service.New(
		service.WithLogger(c.logger),
	); err != nil {
		return err
	}

	// init api
	c.api = api.New(
		api.WithLogger(c.logger),
		api.WithService(c.srv),
	)

	return nil
}

// SetDynamicConfig load dynamic config
func (c *WebScoketComponent) SetDynamicConfig(nf *platformConf.NodeConfig) error {
	return nil
}

// OnConfigChanged modify config
func (c *WebScoketComponent) OnConfigChanged(*platformConf.NodeConfig) error {
	return micro.ErrNeedRestart
}

// SetupHandler setup handler
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
	}
	return nil
}
