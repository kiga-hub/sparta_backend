package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/kiga-hub/sparta_backend/pkg/utils"
)

func main() {
	cmd := exec.Command("python", "./vtk_version.py")
	cmd.Dir = "/home/workspace/project/sparta_backend/script"

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(output), err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		log.Fatal(string(output), err)
	}

	fmt.Println(string(output))
}
