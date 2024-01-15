package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("python", "./vtk_version.py")
	cmd.Dir = "/home/workspace/project/websocket/script"

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(output), err)
	}

	fmt.Println(string(output))
}
