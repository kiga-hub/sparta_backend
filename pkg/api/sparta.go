package api

import (
	"errors"
	"net/http"

	"github.com/kiga-hub/websocket/pkg/models"
	"github.com/kiga-hub/websocket/pkg/utils"
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

	g.POST("/import", s.importSTL).
		SetOperationId("importSTL").
		SetSummary("导入ASCII码形式STL几何").
		SetDescription("导入ASCII码形式STL几何").
		AddParamFile("file", "file", true).
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

func (s *Server) importSTL(c echo.Context) error {
	// handle upload file
	uploadDir, err := s.srv.HandleUploadFile(c)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, errors.New("导入模型失败")))
	}

	// parse import file
	exportInfo, err := s.srv.ParseImportFile(uploadDir)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, err))
	}

	return c.JSON(http.StatusOK, utils.SuccessJSONData(string(exportInfo)))
}
