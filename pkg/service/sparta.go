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

// CreatingParticles -
func (s *Service) CreatingParticles(sparta *models.Sparta) interface{} {

	// process parameters
	sparta.ProcessSparta(GetConfig().DataDir)
	s.logger.Info("Sparta")

	// calculate particles. it will calculate the in.circle file .and generate the out (**)
	s.CalculateSpartaResult()

	return "OK"
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
func (s *Service) ParseImportFile(exportFile string) (*models.SpartaResultDirectory, error) {
	if !utils.IsFileExist(exportFile) {
		return nil, errors.New("import file is not exist")
	}

	file, err := os.Open(exportFile)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	defer file.Close()

	fmt.Println("file.dir: ", filepath.Dir(exportFile))
	fmt.Println("file.name: ", filepath.Base(exportFile))

	stlName := exportFile
	fmt.Println("stlName: ", stlName)

	// generate surf file directory
	surfName := filepath.Join(GetConfig().DataDir, utils.UploadDirName, strings.TrimSuffix(filepath.Base(exportFile), filepath.Ext(filepath.Base(exportFile)))+".surf")
	fmt.Println("surfName: ", surfName)

	// convert to surf file
	{
		cmd := exec.Command("pvpython", "stl2surf.py", stlName, surfName)
		cmd.Dir = GetConfig().ExecDir // "/home/sparta/tools"

		// // read file content
		// data, err := io.ReadAll(file)
		// if err != nil {
		// 	s.logger.Error(err)
		// 	return nil, err
		// }

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}

		// Start executing the command
		if err := cmd.Start(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}

		// Read the command's output in a separate goroutine to prevent blocking
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}

		// Print the output
		fmt.Printf("The output: %s\n", output)
		fmt.Printf("%s\n", output)
	}

	// exec python script. calculate the in. file
	{
		cmd := exec.Command(GetConfig().SpaExec)
		cmd.Dir = GetConfig().DataDir

		file, err := os.Open("/home/workspace/project/sparta_backend/data/in.circle")
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}
		defer file.Close()

		cmd.Stdin = file

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}

		// Start executing the command
		if err := cmd.Start(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return nil, err
		}

		// Read the command's output in a separate goroutine to prevent blocking
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
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

func (s *Service) CalculateSpartaResult() {
	cmd := exec.Command(GetConfig().SpaExec)
	cmd.Dir = GetConfig().DataDir
	// Open the file
	file, err := os.Open("/home/workspace/project/sparta_backend/data/in.circle")
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}
	defer file.Close()
	// Redirect the command's stdin to the file
	cmd.Stdin = file

	// redirect the command's stdin to the string
	// cmd.Stdin = strings.NewReader(dataStr)

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

	// Print the output
	fmt.Printf("The output: %s\n", output)
	fmt.Printf("%s\n", output)

	// Format the output
	result := fmt.Sprintf("%s", output)

	// Write the result to the client
	fmt.Println(result)
}
