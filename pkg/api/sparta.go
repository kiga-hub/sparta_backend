package api

import (
	"errors"
	"net/http"

	"github.com/kiga-hub/sparta_backend/pkg/models"
	"github.com/kiga-hub/sparta_backend/pkg/service"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// setupRESTfulApi setupRESTfulApi
func (s *Server) setupRESTfulApi(root echoswagger.ApiRoot, base string) {
	g := root.Group("RESTful API", base+"/sparta")
	g.POST("", s.creatingParticles).
		SetOperationId(`creating particles`).
		SetSummary("creating particles").
		AddParamBody(models.Sparta{}, "body", "sparta structure", true).
		SetDescription(`creating particles`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.POST("/import", s.importSTL).
		SetOperationId("importSTL").
		SetSummary("Import STL geometry in ASCII format").
		SetDescription("Import STL geometry in ASCII format").
		AddParamFile("file", "file", true).
		AddParamQuery("", "type", "name", true).
		AddResponse(http.StatusOK, ``, nil, nil)
}

// creatingParticles -
func (s *Server) creatingParticles(c echo.Context) error {
	// parse c to models.Sparta
	var sparta models.Sparta
	sparta.MixtureType = make(map[string]string)

	if err := c.Bind(&sparta); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	result := s.srv.ConvertToParaview(&sparta)

	return c.JSON(http.StatusOK, result)
}

// importSTL -
func (s *Server) importSTL(c echo.Context) error {
	// handle upload file. a.stl
	stlDir, err := s.srv.HandleUploadFile(c)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, errors.New("Failure to import model!")))
	}

	fileType := c.QueryParam("type")
	if fileType == "vss" {
		fileName, err := utils.GetFileName(stlDir)
		if err != nil {
			return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, err))
		}
		service.VssFileName = fileName

	}

	if fileType == "species" {
		fileName, err := utils.GetFileName(stlDir)
		if err != nil {
			return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, err))
		}
		service.SpeciesFileName = fileName
	}

	if fileType == "stl" {
		// parse import file. convert stl to surf
		result, err := s.srv.ParseImportFile(stlDir)
		if err != nil {
			s.logger.Error(err)
			return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrImportExportCode, utils.ErrImportExportMsg, err))
		}
		return c.JSON(http.StatusOK, result)
	}

	return c.JSON(http.StatusOK, "OK")
}
