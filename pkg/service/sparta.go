package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kiga-hub/sparta_backend/pkg/models"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// ConvertToParaview -
func (s *Service) ConvertToParaview(sparta *models.Sparta) interface{} {

	// bind file name. be used in other functions. e.g. generate in.circle file.
	sparta.BindFileType(models.SpeciesFileType, models.SpeciesFileName)
	sparta.BindFileType(models.VssFileType, models.VssFileName)

	// if models.GlobalSurfName == "" {
	// 	models.GlobalSurfName = GetConfig().DataDir
	// }
	surfName := strings.Replace(sparta.UploadStlName, "stl", "surf", -1)

	// process parameters
	circleName, err := sparta.ProcessSparta(GetConfig().DataDir, surfName)
	if err != nil {
		return err
	}

	if err := sparta.EditVSSFile(filepath.Join(GetConfig().DataDir, "co2.vss")); err != nil {
		return err
	}

	// {
	// 	// cmpute particles. it will compute the in.circle file .and generate the out (**)
	// 	// s.ComputeSpartaResult(circleName)
	// 	go models.ComputeSpartaResult(circleName, GetConfig().SpaExec)

	// 	// convert to paraview file
	// 	if sparta.IsDumpGrid {
	// 		models.Grid2Paraview(filepath.Dir(circleName), GetConfig().ScriptDir)
	// 		// s.Grid2Paraview(filepath.Dir(circleName))
	// 	}
	// }
	return "ok: " + circleName
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

	// get current time. and convert to 20060102150405
	// currentTime := time.Now().Format("20060102150405")

	// imnpiort file path
	stlDir := filepath.Join(GetConfig().DataDir, file.Filename) // filepath.Join(GetConfig().DataDir, currentTime, file.Filename)

	// create upload dir if not exist
	if err := utils.MakeDirIfNotExist(stlDir); err != nil {
		s.logger.Error(err)
		return "", err
	}

	f, err := os.Create(stlDir)
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

	return stlDir, nil
}

// ParseImportFile -
func (s *Service) ParseImportFile(stlFile string) (*models.SpartaResultDirectory, error) {
	if !utils.IsFileExist(stlFile) {
		return nil, errors.New("import file is not exist")
	}

	file, err := os.Open(stlFile)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer file.Close()

	stlName := stlFile
	surfName := filepath.Join(filepath.Dir(stlName), strings.Replace(filepath.Base(stlFile), filepath.Ext(stlFile), ".surf", -1))

	// convert to surf file
	{
		cmd := exec.Command("pvpython", "stl2surf.py", stlName, surfName)
		cmd.Dir = GetConfig().ScriptDir

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf(utils.ErrorMsg+"ParseImportFile cmd.StdoutPipe()", err)
			return nil, err
		}

		// Start executing the command
		if err := cmd.Start(); err != nil {
			fmt.Printf(utils.ErrorMsg+"ParseImportFile  cmd.Start();", err)
			return nil, err
		}

		// Read the command's output in a separate goroutine to prevent blocking
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(utils.ErrorMsg+"ParseImportFile io.ReadAll(stdout)", err)
			return nil, err
		}
		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf(utils.ErrorMsg+"ParseImportFile cmd.Wait()", err)
			return nil, err
		}

		// Print the output
		fmt.Printf("The output: %s\n", output)
		fmt.Printf("%s\n", output)
	}

	resultInfo := &models.SpartaResultDirectory{
		StlDir:  stlName,
		SurfDir: surfName,
	}

	return resultInfo, nil
}

// ComputeSpartaResult -
func (s *Service) ComputeSpartaResult(circleName string) string {
	cmd := exec.Command(GetConfig().SpaExec)
	cmd.Dir = filepath.Dir(circleName)
	// do spar_ < in.circle
	file, err := os.Open(circleName)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}
	defer file.Close()

	// Redirect the command's stdin to the file
	cmd.Stdin = file

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	// Start executing the command
	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	// Read the command's output in a separate goroutine to prevent blocking
	output, err := io.ReadAll(stdout)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	// Print the output
	fmt.Printf("The output: %s\n", output)
	fmt.Printf("%s\n", output)

	// Format the output
	result := string(output)

	// Write the result to the client
	fmt.Println(result)
	return filepath.Dir(circleName)
}

// Grid2Paraview -
func (s *Service) Grid2Paraview(dir string) {
	go func() {
		// do grid2paraview. pvpython grid2paraview.py circle.txt output -r tmp.grid.1000
		txtFile := filepath.Join(dir, "in.txt")
		outputDir := dir + "/output/"
		tmpGridDir := filepath.Join(dir, "tmp.grid.*")

		// Delete the outputDir directory, TODO need to keep historical files
		if err := utils.ClearDir(outputDir); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		cmd := exec.Command("pvpython", "grid2paraview.py", txtFile, outputDir, "-r", tmpGridDir)
		cmd.Dir = filepath.Join(GetConfig().ScriptDir, "paraview")

		// Create a pipe to capture the command's output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		// Start executing the command
		if err := cmd.Start(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		// Read the command's output in a separate goroutine to prevent blocking
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}
		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		// Print the output
		fmt.Printf("The output: %s\n", output)
		fmt.Printf("%s\n", output)

		// Format the output
		result := fmt.Sprintf("%s", output)

		// Write the result to the client
		fmt.Println(result)
	}()
}
