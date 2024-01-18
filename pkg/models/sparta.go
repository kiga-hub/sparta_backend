package models

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/kiga-hub/sparta_backend/pkg/utils"
)

// Sparta -
type Sparta struct {
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

// ProcessSparta -
func (c *Sparta) ProcessSparta() interface{} {
	// Convert the Data field to a string with newline-separated key-value pairs
	var dataStrs []string
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable x index", c.Data["variable x index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable y index", c.Data["variable y index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable z index", c.Data["variable z index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable n equal", c.Data["variable n equal"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable fnum equal", c.Data["variable fnum equal"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "seed", c.Data["seed"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "dimension", c.Data["dimension"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global nrho", c.Data["global nrho"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global fnum", c.Data["global fnum"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "timestep", c.Data["timestep"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global gridcut", c.Data["global gridcut"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global surfmax", c.Data["global surfmax"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "boundary", c.Data["boundary"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_box", c.Data["create_box"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_grid", c.Data["create_grid"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "balance_grid", c.Data["balance_grid"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "species ar.species", c.Data["species ar.species"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar frac", c.Data["mixture air Ar frac"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air group", c.Data["mixture air group"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar vstream", c.Data["mixture air Ar vstream"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix in emit/face air", c.Data["fix in emit/face air"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide vss air", c.Data["collide vss air"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "read_surf", c.Data["read_surf"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_collide 1 diffuse", c.Data["surf_collide 1 diffuse"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_modify", c.Data["surf_modify"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_particles air n", c.Data["create_particles air n"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix", c.Data["fix"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 yes"))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute g grid all all", c.Data["compute g grid all all"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute max reduce max", c.Data["compute max reduce max"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats_style", c.Data["stats_style"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats", c.Data["stats"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", c.Data["run"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 no"))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", c.Data["run"]))

	dataStr := strings.Join(dataStrs, "\n")

	// Print the result
	fmt.Println(dataStr)

	cmd := exec.Command("/home/spa_")
	cmd.Dir = "/home/sparta-13Apr2023/bench"
	// Open the file
	// file, err := os.Open("/home/sparta-13Apr2023/bench/in.sphere")
	// if err != nil {
	// 	fmt.Printf(utils.ErrorMsg, err)
	// 	return false
	// }
	// defer file.Close()

	// Redirect the command's stdin to the file
	// cmd.Stdin = file

	// redirect the command's stdin to the string
	cmd.Stdin = strings.NewReader(dataStr)

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return err
	}

	// Start executing the command
	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return err
	}

	// Read the command's output in a separate goroutine to prevent blocking
	output, err := io.ReadAll(stdout)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return err
	}

	// Print the output
	fmt.Printf("The output: %s\n", output)
	fmt.Printf("%s\n", output)

	// Format the output
	result := fmt.Sprintf("%s", output)

	return result
}
