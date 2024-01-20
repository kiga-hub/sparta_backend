package models

import (
	"fmt"
	"os"
)

// // Sparta -
// type Sparta struct {
// 	Message string            `json:"message"`
// 	Data    map[string]string `json:"data"`
// }

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
	GlobalNrho      string   `json:"global_nrho"`       // 来流分子数密度
	GlobalFnum      string   `json:"global_fnum"`       // 真实与模拟分子数之比
	SurfCollideType string   `json:"surf_collide_type"` // 粒子与物面碰撞模型 diffuse,specular
	CollideAlpha    string   `json:"collide_alpha"`     // 粒子与粒子碰撞模型 硬:1,软:1.4
	WallTemperature string   `json:"wall_temperature"`  // 壁面温度
	Reflectivity    string   `json:"reflectivity"`      // 反射率
	MixtureType     []string `json:"mixture_type"`      // 混合物 N2 CO2 O2
	Temperature     string   `json:"temperature"`       // 温度
	VStreamX        string   `json:"v_stream_x"`        // 来流速度 x
	VStreamY        string   `json:"v_stream_y"`        // 来流速度 y
	VStreamZ        string   `json:"v_stream_z"`        // 来流速度 z

	// 计算速度
	ComputeSpeed []string `json:"compute_speed"` // 计算网格速度u\ v\ w

	// 计算热力学
	ComputeThermo []string `json:"compute_thermo"` // 计算热力学 temp \ press

	// 计算热流密度
	ComputeHeat []string `json:"compute_heat"` // 计算热流密度 heat_x\ heat_y\ heat_z

	TimeStep string `json:"time_step"` // 时间步长
	Run      string `json:"run"`       // 计算步数
}

func (c *Sparta) ProcessSparta() {
	fmt.Println("Process Sparta: ", c)
	// 打开文件用于写入
	file, err := os.Create("./data/in.circle")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var pre = `################################################################################
# 2d flow around a circle
#
# Note:
#  - The "comm/sort” option to the “global” command is used to match MPI runs.
#  - The “twopass” option is used to match Kokkos runs.
# The "comm/sort" and "twopass" options should not be used for production runs.
################################################################################
`
	fmt.Fprintf(file, pre)
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "seed			 %s\n", "12345")
	// 根据 Sparta 结构体字段写入文件
	fmt.Fprintf(file, "dimension        %s\n", c.Dimension)
	fmt.Fprintf(file, "global           gridcut %s comm/sort %s\n", "0.0", "yes")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "boundary         %s%s %s%s %s%s\n", c.BoundaryXLO, c.BoundaryXHI, c.BoundaryYLO, c.BoundaryYHI, c.BoundaryZLO, c.BoundaryZHI)
	fmt.Fprintf(file, "create_box       %s %s %s %s %s %s\n", c.CreateBoxXMin, c.CreateBoxXMax, c.CreateBoxYMin, c.CreateBoxYMax, c.CreateBoxZMin, c.CreateBoxZMax)
	fmt.Fprintf(file, "create_grid      %s %s %s\n", c.CreateGridX, c.CreateGridY, c.CreateGridZ)
	fmt.Fprintf(file, "balance_grid     %s %s\n", "rcb", "cell")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "global           nrho %s fnum %s\n", c.GlobalNrho, c.GlobalFnum)
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "species          %s %s\n", "co2.species", "CO2")

	// 解析MixtureType 字段
	var mixtureType string
	for _, v := range c.MixtureType {
		mixtureType += v + " "
	}

	fmt.Fprintf(file, "mixture          %s %s %s %s %s %s %s %s\n", "air", mixtureType, "vstream", c.VStreamX, c.VStreamY, c.VStreamZ, "temp", c.Temperature)
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "read_surf        %s %s %s\n", "b.surf", "scale", "0.001 0.001 0.001")
	fmt.Fprintf(file, "surf_collide     %s %s %s %s\n", "1", c.SurfCollideType, c.WallTemperature, c.Reflectivity)
	fmt.Fprintf(file, "surf_modify      %s %s %s\n", "all", "collide", "1")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "collide          %s %s %s\n", "vss", "air", "co2.vss")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "fix              %s %s %s %s %s\n", "in", "emit/face", "air", "xlo", "twopass")

	// 解析ComputeSpeed 字段
	var computeSpeed string
	for _, v := range c.ComputeSpeed {
		computeSpeed += v + " "
	}
	fmt.Fprintf(file, "compute          %s %s %s %s %s\n", "1", "grid", "all", "species", computeSpeed)
	fmt.Fprintf(file, "fix              %s %s %s %s %s %s\n", "1", "ave/grid", "all", "10", "100", "1000 c_1[*]")
	fmt.Fprintf(file, "\n")

	// 解析ComputeHeat 字段
	var computeHeat string
	for _, v := range c.ComputeHeat {
		computeHeat += v + " "
	}
	fmt.Fprintf(file, "compute          %s %s %s %s %s\n", "2", "eflux/grid", "all", "species", computeHeat)
	fmt.Fprintf(file, "fix              %s %s %s %s %s %s\n", "2", "ave/grid", "all", "10", "100", "1000 c_2[*]")
	fmt.Fprintf(file, "\n")

	// 解析ComputeThermo 字段
	var computeThermo string
	for _, v := range c.ComputeThermo {
		computeThermo += v + " "
	}

	fmt.Fprintf(file, "compute          %s %s %s %s %s\n", "3", "thermal/grid", "all", "species", computeThermo)
	fmt.Fprintf(file, "fix              %s %s %s %s %s %s\n", "3", "ave/grid", "all", "10", "100", "1000 c_3[*]")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "dump             %s %s %s %s %s %s %s %s %s %s %s %s\n", "1", "grid", "all", "1000", "tmp.grid.*", "id", "xc", "yc", "zc", "f_1[*]", "f_2[*]", "f_3[*]")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "write_grid       %s %s\n", "data.grid", "")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "timestep         %s\n", c.TimeStep)
	fmt.Fprintf(file, "\n")

	var dump = `
dump                2 image all 100 image.*.ppm type type pdiam 0.001 &
			surf proc 0.01 size 1024 1024 zoom 1.75 &
			gline no 0.005
dump_modify	    2 pad 4
`

	fmt.Fprintf(file, dump)
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "stats            %s\n", "100")
	fmt.Fprintf(file, "stats_style      %s %s %s %s %s %s\n", "step", "cpu", "np", "nattempt", "ncoll", "nscoll nscheck")
	fmt.Fprintf(file, "\n")

	fmt.Fprintf(file, "run              %s\n", c.Run)
	fmt.Fprintf(file, "\n")

	// Done
	fmt.Println("Done")

}

// // ProcessSparta -
// func (c *Sparta) ProcessSparta() interface{} {
// 	// Convert the Data field to a string with newline-separated key-value pairs
// 	var dataStrs []string
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable x index", c.Data["variable x index"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable y index", c.Data["variable y index"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable z index", c.Data["variable z index"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable n equal", c.Data["variable n equal"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable fnum equal", c.Data["variable fnum equal"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "seed", c.Data["seed"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "dimension", c.Data["dimension"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global nrho", c.Data["global nrho"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global fnum", c.Data["global fnum"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "timestep", c.Data["timestep"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global gridcut", c.Data["global gridcut"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global surfmax", c.Data["global surfmax"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "boundary", c.Data["boundary"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_box", c.Data["create_box"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_grid", c.Data["create_grid"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "balance_grid", c.Data["balance_grid"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "species ar.species", c.Data["species ar.species"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar frac", c.Data["mixture air Ar frac"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air group", c.Data["mixture air group"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar vstream", c.Data["mixture air Ar vstream"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix in emit/face air", c.Data["fix in emit/face air"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide vss air", c.Data["collide vss air"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "read_surf", c.Data["read_surf"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_collide 1 diffuse", c.Data["surf_collide 1 diffuse"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_modify", c.Data["surf_modify"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_particles air n", c.Data["create_particles air n"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix", c.Data["fix"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 yes"))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute g grid all all", c.Data["compute g grid all all"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute max reduce max", c.Data["compute max reduce max"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats_style", c.Data["stats_style"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats", c.Data["stats"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", c.Data["run"]))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 no"))
// 	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", c.Data["run"]))

// 	dataStr := strings.Join(dataStrs, "\n")

// 	// Print the result
// 	fmt.Println(dataStr)

// 	cmd := exec.Command("/home/spa_")
// 	cmd.Dir = "/home/sparta-13Apr2023/bench"
// 	// Open the file
// 	// file, err := os.Open("/home/sparta-13Apr2023/bench/in.sphere")
// 	// if err != nil {
// 	// 	fmt.Printf(utils.ErrorMsg, err)
// 	// 	return false
// 	// }
// 	// defer file.Close()

// 	// Redirect the command's stdin to the file
// 	// cmd.Stdin = file

// 	// redirect the command's stdin to the string
// 	cmd.Stdin = strings.NewReader(dataStr)

// 	// Create a pipe to capture the command's output
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		fmt.Printf(utils.ErrorMsg, err)
// 		return err
// 	}

// 	// Start executing the command
// 	if err := cmd.Start(); err != nil {
// 		fmt.Printf(utils.ErrorMsg, err)
// 		return err
// 	}

// 	// Read the command's output in a separate goroutine to prevent blocking
// 	output, err := io.ReadAll(stdout)
// 	if err != nil {
// 		fmt.Printf(utils.ErrorMsg, err)
// 		return err
// 	}

// 	// Print the output
// 	fmt.Printf("The output: %s\n", output)
// 	fmt.Printf("%s\n", output)

// 	// Format the output
// 	result := fmt.Sprintf("%s", output)

// 	return result
// }
