package service

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kiga-hub/websocket/pkg/models"
	"github.com/kiga-hub/websocket/pkg/utils"
	"github.com/labstack/echo/v4"
)

// CreatingParticles -
func (s *Service) CreatingParticles(sparta *models.Sparta) interface{} {

	result := sparta.ProcessSparta()
	s.logger.Info("Sparta: ", result)

	return result
}

// HandleUploadFile handle upload file
func (s *Service) HandleUploadFile(c echo.Context) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}

	// open the file
	src, err := file.Open()
	if err != nil {
		s.logger.Error("open file failed")
		return "", c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrGetDataCode, utils.ErrGetDataMsg, err))
	}
	defer src.Close()

	// clear upload dir
	// emptyDIr := filepath.Join(GetConfig().DataDir, utils.UploadDirName)
	// if err := utils.ClearDir(emptyDIr); err != nil {
	// 	s.logger.Error(err)
	// 	return "", err
	// }

	// imnpiort file path
	uploadDir := filepath.Join(GetConfig().DataDir, utils.UploadDirName, file.Filename)

	// create upload dir if not exist
	if err := utils.MakeDirIfNotExist(uploadDir); err != nil {
		s.logger.Error(err)
		return "", err
	}

	f, err := os.Create(uploadDir)
	if err != nil {
		s.logger.Error(err)
		return "", err
	}
	defer f.Close()

	// copy the file to the destination
	if _, err = io.Copy(f, src); err != nil {
		s.logger.Error(err)
		return "", err
	}

	return uploadDir, nil
}

// ParseImportFile -
func (s *Service) ParseImportFile(exportFile string) ([]byte, error) {
	if !utils.IsFileExist(exportFile) {
		return nil, errors.New("import file is not exist")
	}

	file, err := os.Open(exportFile)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer file.Close()

	// read file content
	data, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return data, nil
}
