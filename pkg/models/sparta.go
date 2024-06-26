package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kiga-hub/sparta_backend/pkg/utils"
)

// var GlobalSurfName string

// SpartaResultDirectory -
type SpartaResultDirectory struct {
	StlDir  string `json:"stl_dir"`
	SurfDir string `json:"surf_dir"`
	OutDir  string `json:"out_dir"`
}

const (
	// VssFileType -
	VssFileType = "vss"
	// SpeciesFileType -
	SpeciesFileType = "species"
)

var (
	// SpeciesFileName co2.species
	SpeciesFileName = "co2.species"
	// VssFileName co2.vss
	VssFileName = "co2.vss"
)

// Sparta -
type Sparta struct {
	Dimension         string            `json:"dimension"`
	BoundaryXLO       string            `json:"boundary_xlo"` // o a r s
	BoundaryXHI       string            `json:"boundary_xhi"`
	BoundaryYLO       string            `json:"boundary_ylo"`
	BoundaryYHI       string            `json:"boundary_yhi"`
	BoundaryZLO       string            `json:"boundary_zlo"`
	BoundaryZHI       string            `json:"boundary_zhi"`
	CreateBoxXMin     string            `json:"create_box_x_min"`
	CreateBoxXMax     string            `json:"create_box_x_max"`
	CreateBoxYMin     string            `json:"create_box_y_min"`
	CreateBoxYMax     string            `json:"create_box_y_max"`
	CreateBoxZMin     string            `json:"create_box_z_min"`
	CreateBoxZMax     string            `json:"create_box_z_max"`
	CreateGridX       string            `json:"create_grid_x"`
	CreateGridY       string            `json:"create_grid_y"`
	CreateGridZ       string            `json:"create_grid_z"`
	GlobalNrho        string            `json:"global_nrho"`
	GlobalFnum        string            `json:"global_fnum"`
	SurfCollideType   string            `json:"surf_collide_type"` // diffuse,specular
	CollideAlpha      string            `json:"collide_alpha"`     // hard:1,soft:1.4
	WallTemperature   string            `json:"wall_temperature"`
	Reflectivity      string            `json:"reflectivity"`
	MixtureType       map[string]string `json:"mixture_type"`     // N2 CO2 O2
	MixtureTypeStr    string            `json:"mixture_type_str"` // N2 CO2 O2
	Temperature       string            `json:"temperature"`
	VStreamX          string            `json:"v_stream_x"`     // vstream x
	VStreamY          string            `json:"v_stream_y"`     // vstream y
	VStreamZ          string            `json:"v_stream_z"`     // vstream z
	ComputeSpeed      []string          `json:"compute_speed"`  //  speed u\ v\ w
	ComputeThermo     []string          `json:"compute_thermo"` // Thermo temp \ press
	ComputeHeat       []string          `json:"compute_heat"`   // Heat heat_x\ heat_y\ heat_z
	TimeStep          string            `json:"time_step"`
	Run               string            `json:"run"`
	UploadStlName     string            `json:"upload_stl_name"`
	IsDumpGrid        bool              `json:"is_dump_grid"`         // Output result document
	IsGridToParaView  bool              `json:"is_grid_to_paraview"`  // Convert to visual format
	DumpGridNumber    string            `json:"dump_grid_number"`     // Used to dump grid
	DumpComputeSpeed  []string          `json:"dump_compute_speed"`   // speed u\ v\ w
	DumpComputeThermo []string          `json:"dump_compute_thermo"`  //temp \ press
	IsDumpComputeHeat bool              `json:"is_dump_compute_heat"` // Heat heat_x\ heat_y\ heat_z
	IsGridCoordinate  bool              `json:"is_grid_coordinate"`   // Grid coordinate
	SpeciesFileName   string            `json:"-"`
	VssFileName       string            `json:"-"`
}

// BindFileType -
func (s *Sparta) BindFileType(fileType, fileName string) {
	if fileType == VssFileType {
		s.VssFileName = fileName
	}

	if fileType == SpeciesFileType {
		s.SpeciesFileName = fileName

	}
}

// SetFileType -
func (s *Sparta) SetFileType() {
	s.SpeciesFileName = SpeciesFileName
	s.VssFileName = VssFileName
}

// ProcessSparta -
func (s *Sparta) ProcessSparta(dir, surfName string) (string, error) {
	// fmt.Println("Process Sparta: ", c)
	s.SetFileType()
	// create circle txt file
	txtFile, err := os.Create(filepath.Join(dir, "in.txt"))
	if err != nil {
		fmt.Println("os.Creat in.txt err", err)
		return "", err
	}
	defer txtFile.Close()

	fmt.Fprintf(txtFile, "\n")
	fmt.Fprintf(txtFile, "dimension        %s\n", s.Dimension)
	fmt.Fprintf(txtFile, "\n")
	fmt.Fprintf(txtFile, "create_box       %s %s %s %s %s %s\n", s.CreateBoxXMin, s.CreateBoxXMax, s.CreateBoxYMin, s.CreateBoxYMax, s.CreateBoxZMin, s.CreateBoxZMax)
	fmt.Fprintf(txtFile, "read_grid        %s \n", dir+"/data.grid")
	fmt.Fprintf(txtFile, "\n")

	// write to in.circle file
	inFile, err := os.Create(filepath.Join(dir, "in.circle"))
	if err != nil {
		fmt.Println("os.Creat in.circle err", err)
		return "", err
	}
	defer inFile.Close()

	var pre = `################################################################################
# 2d flow around a circle
#
# Note:
#  - The "comm/sort” option to the “global” command is used to match MPI runs.
#  - The “twopass” option is used to match Kokkos runs.
# The "comm/sort" and "twopass" options should not be used for production runs.
################################################################################
`
	fmt.Fprintf(inFile, pre)
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "seed			 %s\n", "12345")
	// Write to a file according to the Sparta struct fields
	fmt.Fprintf(inFile, "dimension        %s\n", s.Dimension)
	fmt.Fprintf(inFile, "global           gridcut %s comm/sort %s\n", "0.0", "yes")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "boundary         %s%s %s%s %s%s\n", s.BoundaryXLO, s.BoundaryXHI, s.BoundaryYLO, s.BoundaryYHI, s.BoundaryZLO, s.BoundaryZHI)
	fmt.Fprintf(inFile, "create_box       %s %s %s %s %s %s\n", s.CreateBoxXMin, s.CreateBoxXMax, s.CreateBoxYMin, s.CreateBoxYMax, s.CreateBoxZMin, s.CreateBoxZMax)
	fmt.Fprintf(inFile, "create_grid      %s %s %s\n", s.CreateGridX, s.CreateGridY, s.CreateGridZ)
	fmt.Fprintf(inFile, "balance_grid     %s %s\n", "rcb", "cell")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "global           nrho %s fnum %s\n", s.GlobalNrho, s.GlobalFnum)
	fmt.Fprintf(inFile, "\n")

	// using json ummarshal convert s.MixtureTypeStr to s.MixtureType（map[string）string
	s.MixtureType = make(map[string]string)
	if s.MixtureTypeStr != "" {
		err := json.Unmarshal([]byte(s.MixtureTypeStr), &s.MixtureType)
		if err != nil {
			fmt.Println("json.Unmarshal err", err)
			return "", err
		}
	}

	// parse MixtureType
	var mixtureType string
	for key := range s.MixtureType {
		mixtureType += key + " "
	}

	fmt.Fprintf(inFile, "species          %s %s\n", s.SpeciesFileName, mixtureType)

	/*
		species air.species N2 CO2 O2
		mixture air N2 vstream 100.0 0 0 temp frac 0.2
		mixture air CO2 vstream 100.0 0 0 temp frac 0.7
		mixture air O2 vstream 100.0 0 0 temp frac 0.1
	*/
	for key, value := range s.MixtureType { //value
		fmt.Fprintf(inFile, "mixture          %s %s %s %s %s %s %s %s %s %s\n", "air", key, "vstream", s.VStreamX, s.VStreamY, s.VStreamZ, "temp", s.Temperature, "frac", value)
		// fmt.Fprintf(inFile, "mixture          %s %s %s %s %s %s %s %s \n", "air", key, "vstream", s.VStreamX, s.VStreamY, s.VStreamZ, "temp", s.Temperature)
		// mixture air N frac 0.8
		// fmt.Fprintf(inFile, "mixture          %s %s %s %s\n", "air", key, "frac", value)
	}

	// fmt.Fprintf(inFile, "mixture          %s %s %s %s %s %s %s %s\n", "air", mixtureType, "vstream", s.VStreamX, s.VStreamY, s.VStreamZ, "temp", s.Temperature)
	fmt.Fprintf(inFile, "\n")

	// fmt.Fprintf(inFile, "read_surf        %s %s %s\n", filepath.Base(GlobalSurfName), "scale", "0.001 0.001 0.001")
	fmt.Fprintf(inFile, "read_surf        %s %s %s\n", surfName, "scale", "0.001 0.001 0.001") // surfName "b.surf"
	fmt.Fprintf(inFile, "surf_collide     %s %s %s %s\n", "1", s.SurfCollideType, s.WallTemperature, s.Reflectivity)
	fmt.Fprintf(inFile, "surf_modify      %s %s %s\n", "all", "collide", "1")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "collide          %s %s %s\n", "vss", "air", s.VssFileName)
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "fix              %s %s %s %s %s\n", "in", "emit/face", "air", "xlo", "twopass")

	// parse ComputeSpeed
	var computeSpeed string
	for _, v := range s.ComputeSpeed {
		computeSpeed += v + " "
	}
	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "1", "grid", "all", "species", computeSpeed)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "1", "ave/grid", "all", "10", "100", "1000 c_1[*]")
	fmt.Fprintf(inFile, "\n")

	// parse ComputeHeat
	var computeHeat string
	for _, v := range s.ComputeHeat {
		computeHeat += v + " "
	}
	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "2", "eflux/grid", "all", "species", computeHeat)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "2", "ave/grid", "all", "10", "100", "1000 c_2[*]")
	fmt.Fprintf(inFile, "\n")

	// parse ComputeThermo
	var computeThermo string
	for _, v := range s.ComputeThermo {
		computeThermo += v + " "
	}

	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "3", "thermal/grid", "all", "species", computeThermo)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "3", "ave/grid", "all", "10", "100", "1000 c_3[*]")
	fmt.Fprintf(inFile, "\n")

	if s.IsDumpGrid {
		//  dump 1 grid 1000 tmp.grid*.id xc yc zc f_1[*] f_2[*] f_3[*]
		// dumpGridString := s.DumpGridNumber + " grid all " + s.Run + " tmp.grid.* id"
		dumpGridString := s.DumpGridNumber + " grid all " + "1000" + " tmp.grid.* id"

		if s.IsGridCoordinate {
			dumpGridString += " xc yc zc"
		}

		if len(s.DumpComputeSpeed) == 3 {
			dumpGridString += " f_1[*]"
		} else {
			speedMap := map[string]string{"u": "1", "v": "2", "w": "3"}
			for _, speed := range s.DumpComputeSpeed {
				if val, ok := speedMap[speed]; ok {
					dumpGridString += " f_1[" + val + "]"
				}
			}
		}

		if s.IsDumpComputeHeat {
			dumpGridString += " f_2[*]"
		}

		if len(s.DumpComputeThermo) == 2 {
			dumpGridString += " f_3[*]"
		} else {
			thermoMap := map[string]string{"temp": "1", "press": "2"}
			for _, thermo := range s.DumpComputeThermo {
				if val, ok := thermoMap[thermo]; ok {
					dumpGridString += " f_3[" + val + "]"
				}
			}
		}
		fmt.Fprintf(inFile, "dump             %s\n", dumpGridString)
		//fmt.Fprintf(inFile, "dump             %s %s %s %s %s %s %s %s %s %s %s %s\n", "1", "grid", "all", "1000", "tmp.grid.*", "id", "xc", "yc", "zc", "f_1[*]", "f_2[*]", "f_3[*]")
		fmt.Fprintf(inFile, "\n")
	}

	fmt.Fprintf(inFile, "write_grid       %s %s\n", "data.grid", "")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "timestep         %s\n", s.TimeStep)
	fmt.Fprintf(inFile, "\n")

	var dumpString string
	dumpString += "dump                2 image all 1000 " + surfName + ".*.ppm type type pdiam 0.001 &\n"
	dumpString += "			surf proc 0.01 size 1024 1024 zoom 1.75 &\n"
	dumpString += "			gline no 0.005\n"
	dumpString += "dump_modify	    2 pad 4\n"

	// 	var dump = `
	// dump                2 image all 100 image.*.ppm type type pdiam 0.001 &
	// 			surf proc 0.01 size 1024 1024 zoom 1.75 &
	// 			gline no 0.005
	// dump_modify	    2 pad 4
	// `
	fmt.Fprintf(inFile, "%s", dumpString)

	// fmt.Fprintf(inFile, "%s", dump)
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "stats            %s\n", "100")
	fmt.Fprintf(inFile, "stats_style      %s %s %s %s %s %s\n", "step", "cpu", "np", "nattempt", "ncoll", "nscoll nscheck")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "run              %s\n", s.Run)
	fmt.Fprintf(inFile, "\n")

	return filepath.Join(dir, "in.circle"), nil
}

// EditVSSFile - collideAlpha hard:1,soft:1.4
func (s *Sparta) EditVSSFile(fileName string) error {
	// open file
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Modify each line that does not start with a "#"
	for i, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			values := strings.Fields(line)
			if len(values) > 4 {
				values[4] = s.CollideAlpha
				lines[i] = strings.Join(values, "   ")
			}
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Sparta) ToBytes() []byte {
	return []byte(fmt.Sprintf("%v", s))
}

// ComputeSpartaResult -
func ComputeSpartaResult(circleName string, spaExe string) string {
	cmd := exec.Command(spaExe)
	cmd.Dir = filepath.Dir(circleName)

	file, err := os.Open(circleName)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}
	defer file.Close()

	cmd.Stdin = file
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	return filepath.Dir(circleName)
}

// ComputeSpartaResult3 -
func ComputeSpartaResult3(circleName string, spaExe string) string {
	cmd := exec.Command(spaExe)
	cmd.Dir = filepath.Dir(circleName)

	file, err := os.Open(circleName)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}
	defer file.Close()

	cmd.Stdin = file

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	return filepath.Dir(circleName)
}

// ComputeSpartaResult2 -
func ComputeSpartaResult2(circleName string, spaExe string) string {
	cmd := exec.Command(spaExe)
	cmd.Dir = filepath.Dir(circleName)

	fmt.Printf("cmd:%s, name:%s\n", spaExe, circleName)
	// do spar_ < in.circle
	file, err := os.Open(circleName)
	if err != nil {
		fmt.Printf(utils.ErrorMsg+"os.Open(circleName)", err)
		return err.Error()
	}
	defer file.Close()

	// Redirect the command's stdin to the file
	cmd.Stdin = file

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg+"cmd.StdoutPipe()", err)
		return err.Error()
	}

	// Start executing the command
	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg+"cmd.Start()", err)
		return err.Error()
	}

	// Read the command's output in a separate goroutine to prevent blocking
	output, err := io.ReadAll(stdout)
	if err != nil {
		fmt.Printf(utils.ErrorMsg+"io.ReadAll(stdout)", err)
		return err.Error()
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg+"cmd.Wait()", err)
		return err.Error()
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
func Grid2Paraview(dir, scriptDir string) {
	go func() {
		// do grid2paraview. pvpython grid2paraview.py circle.txt output -r tmp.grid.1000
		txtFile := filepath.Join(dir, "in.txt")
		outputDir := dir + "/output/"
		tmpGridDir := filepath.Join(dir, "tmp.grid.*")

		// Delete the outputDir directory, TODO need to keep historical files
		if err := utils.ClearDir(outputDir); err != nil {
			fmt.Printf(utils.ErrorMsg+"Grid2Paraview utils.ClearDir(outputDir)", err)
			return
		}

		cmd := exec.Command("pvpython", "grid2paraview.py", txtFile, outputDir, "-r", tmpGridDir)
		cmd.Dir = filepath.Join(scriptDir, "paraview")

		// Create a pipe to capture the command's output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf(utils.ErrorMsg+"Grid2Paraview cmd.StdoutPipe()", err)
			return
		}

		// Start executing the command
		if err := cmd.Start(); err != nil {
			fmt.Printf(utils.ErrorMsg+"Grid2Paraview cmd.Start()", err)
			return
		}

		// Read the command's output in a separate goroutine to prevent blocking
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(utils.ErrorMsg+"Grid2Paraview io.ReadAll(stdout)", err)
			return
		}
		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf(utils.ErrorMsg+"Grid2Paraview cmd.Wait()", err)
			return
		}

		// Print the output
		fmt.Printf("The output: %s\n", output)
		fmt.Printf("%s\n", output)

		// Format the output
		result := string(output)

		// Write the result to the client
		fmt.Println(result)
	}()
}
