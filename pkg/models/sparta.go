package models

import (
	"fmt"
	"os"
	"path/filepath"
)

// var GlobalSurfName string

// SpartaResultDirectory -
type SpartaResultDirectory struct {
	StlDir  string `json:"stl_dir"`
	SurfDir string `json:"surf_dir"`
	OutDir  string `json:"out_dir"`
}

// Sparta -
type Sparta struct {
	Dimension       string   `json:"dimension"`
	BoundaryXLO     string   `json:"boundary_xlo"` // o a r s
	BoundaryXHI     string   `json:"boundary_xhi"`
	BoundaryYLO     string   `json:"boundary_ylo"`
	BoundaryYHI     string   `json:"boundary_yhi"`
	BoundaryZLO     string   `json:"boundary_zlo"`
	BoundaryZHI     string   `json:"boundary_zhi"`
	CreateBoxXMin   string   `json:"create_box_x_min"`
	CreateBoxXMax   string   `json:"create_box_x_max"`
	CreateBoxYMin   string   `json:"create_box_y_min"`
	CreateBoxYMax   string   `json:"create_box_y_max"`
	CreateBoxZMin   string   `json:"create_box_z_min"`
	CreateBoxZMax   string   `json:"create_box_z_max"`
	CreateGridX     string   `json:"create_grid_x"`
	CreateGridY     string   `json:"create_grid_y"`
	CreateGridZ     string   `json:"create_grid_z"`
	GlobalNrho      string   `json:"global_nrho"`
	GlobalFnum      string   `json:"global_fnum"`
	SurfCollideType string   `json:"surf_collide_type"` // diffuse,specular
	CollideAlpha    string   `json:"collide_alpha"`     // hard:1,soft:1.4
	WallTemperature string   `json:"wall_temperature"`
	Reflectivity    string   `json:"reflectivity"`
	MixtureType     []string `json:"mixture_type"` // N2 CO2 O2
	Temperature     string   `json:"temperature"`
	VStreamX        string   `json:"v_stream_x"`     // vstream x
	VStreamY        string   `json:"v_stream_y"`     // vstream y
	VStreamZ        string   `json:"v_stream_z"`     // vstream z
	ComputeSpeed    []string `json:"compute_speed"`  //  speed u\ v\ w
	ComputeThermo   []string `json:"compute_thermo"` // Thermo temp \ press
	ComputeHeat     []string `json:"compute_heat"`   // Heat heat_x\ heat_y\ heat_z
	TimeStep        string   `json:"time_step"`
	Run             string   `json:"run"`
	UploadStlName   string   `json:"upload_stl_name"`
	IsDumpGrid      bool     `json:"is_dump_grid"`
}

// ProcessSparta -
func (c *Sparta) ProcessSparta(dir, surfName string) string {
	// fmt.Println("Process Sparta: ", c)

	// create circle txt file
	txtFile, err := os.Create(filepath.Join(dir, "in.txt"))
	if err != nil {
		panic(err)
	}
	defer txtFile.Close()

	fmt.Fprintf(txtFile, "\n")
	fmt.Fprintf(txtFile, "dimension        %s\n", c.Dimension)
	fmt.Fprintf(txtFile, "\n")
	fmt.Fprintf(txtFile, "create_box       %s %s %s %s %s %s\n", c.CreateBoxXMin, c.CreateBoxXMax, c.CreateBoxYMin, c.CreateBoxYMax, c.CreateBoxZMin, c.CreateBoxZMax)
	fmt.Fprintf(txtFile, "read_grid        %s \n", dir+"/data.grid")
	fmt.Fprintf(txtFile, "\n")

	// write to in.circle file
	inFile, err := os.Create(filepath.Join(dir, "in.circle"))
	if err != nil {
		panic(err)
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
	// 根据 Sparta 结构体字段写入文件
	fmt.Fprintf(inFile, "dimension        %s\n", c.Dimension)
	fmt.Fprintf(inFile, "global           gridcut %s comm/sort %s\n", "0.0", "yes")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "boundary         %s%s %s%s %s%s\n", c.BoundaryXLO, c.BoundaryXHI, c.BoundaryYLO, c.BoundaryYHI, c.BoundaryZLO, c.BoundaryZHI)
	fmt.Fprintf(inFile, "create_box       %s %s %s %s %s %s\n", c.CreateBoxXMin, c.CreateBoxXMax, c.CreateBoxYMin, c.CreateBoxYMax, c.CreateBoxZMin, c.CreateBoxZMax)
	fmt.Fprintf(inFile, "create_grid      %s %s %s\n", c.CreateGridX, c.CreateGridY, c.CreateGridZ)
	fmt.Fprintf(inFile, "balance_grid     %s %s\n", "rcb", "cell")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "global           nrho %s fnum %s\n", c.GlobalNrho, c.GlobalFnum)
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "species          %s %s\n", "co2.species", "CO2")

	// parse MixtureType
	var mixtureType string
	for _, v := range c.MixtureType {
		mixtureType += v + " "
	}

	fmt.Fprintf(inFile, "mixture          %s %s %s %s %s %s %s %s\n", "air", mixtureType, "vstream", c.VStreamX, c.VStreamY, c.VStreamZ, "temp", c.Temperature)
	fmt.Fprintf(inFile, "\n")

	// fmt.Fprintf(inFile, "read_surf        %s %s %s\n", filepath.Base(GlobalSurfName), "scale", "0.001 0.001 0.001")
	fmt.Fprintf(inFile, "read_surf        %s %s %s\n", "b.surf", "scale", "0.001 0.001 0.001") // surfName
	fmt.Fprintf(inFile, "surf_collide     %s %s %s %s\n", "1", c.SurfCollideType, c.WallTemperature, c.Reflectivity)
	fmt.Fprintf(inFile, "surf_modify      %s %s %s\n", "all", "collide", "1")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "collide          %s %s %s\n", "vss", "air", "co2.vss")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "fix              %s %s %s %s %s\n", "in", "emit/face", "air", "xlo", "twopass")

	// parse ComputeSpeed
	var computeSpeed string
	for _, v := range c.ComputeSpeed {
		computeSpeed += v + " "
	}
	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "1", "grid", "all", "species", computeSpeed)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "1", "ave/grid", "all", "10", "100", "1000 c_1[*]")
	fmt.Fprintf(inFile, "\n")

	// parse ComputeHeat
	var computeHeat string
	for _, v := range c.ComputeHeat {
		computeHeat += v + " "
	}
	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "2", "eflux/grid", "all", "species", computeHeat)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "2", "ave/grid", "all", "10", "100", "1000 c_2[*]")
	fmt.Fprintf(inFile, "\n")

	// parse ComputeThermo
	var computeThermo string
	for _, v := range c.ComputeThermo {
		computeThermo += v + " "
	}

	fmt.Fprintf(inFile, "compute          %s %s %s %s %s\n", "3", "thermal/grid", "all", "species", computeThermo)
	fmt.Fprintf(inFile, "fix              %s %s %s %s %s %s\n", "3", "ave/grid", "all", "10", "100", "1000 c_3[*]")
	fmt.Fprintf(inFile, "\n")

	if c.IsDumpGrid {
		fmt.Fprintf(inFile, "dump             %s %s %s %s %s %s %s %s %s %s %s %s\n", "1", "grid", "all", "1000", "tmp.grid.*", "id", "xc", "yc", "zc", "f_1[*]", "f_2[*]", "f_3[*]")
		fmt.Fprintf(inFile, "\n")
	}

	fmt.Fprintf(inFile, "write_grid       %s %s\n", "data.grid", "")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "timestep         %s\n", c.TimeStep)
	fmt.Fprintf(inFile, "\n")

	var dump = `
dump                2 image all 100 image.*.ppm type type pdiam 0.001 &
			surf proc 0.01 size 1024 1024 zoom 1.75 &
			gline no 0.005
dump_modify	    2 pad 4
`

	fmt.Fprintf(inFile, dump)
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "stats            %s\n", "100")
	fmt.Fprintf(inFile, "stats_style      %s %s %s %s %s %s\n", "step", "cpu", "np", "nattempt", "ncoll", "nscoll nscheck")
	fmt.Fprintf(inFile, "\n")

	fmt.Fprintf(inFile, "run              %s\n", c.Run)
	fmt.Fprintf(inFile, "\n")

	// fmt.Println("Done")
	return filepath.Join(dir, "in.circle")
}
