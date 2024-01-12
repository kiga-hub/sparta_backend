package api

import (
	"net/http"

	"github.com/kiga-hub/websocket/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// setupRESTfulApi setupRESTfulApi
func (s *Server) setupRESTfulApi(root echoswagger.ApiRoot, base string) {
	g := root.Group("RESTful API", base+"/sparta")
	g.POST("", s.CreatingParticles).
		SetOperationId(`creating particles`).
		SetSummary("creating particles").
		AddParamBody(models.Sparta{}, "body", "sparta structure", true).
		SetDescription(`creating particles`).
		AddResponse(http.StatusOK, ``, nil, nil)
}

// CreatingParticles -
func (s *Server) CreatingParticles(c echo.Context) error {

	// parse c to models.Sparta
	var sparta models.Sparta
	if err := c.Bind(&sparta); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	result := s.srv.CreatingParticles(&sparta)

	return c.JSON(http.StatusOK, result)
}
