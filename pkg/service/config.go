package service

import "github.com/spf13/viper"

const (
	dataDir          = "service.data_dir"
	rootDir          = "service.root_dir"
	spartaPytonTools = "service.sparta_python_tools"
	execDir          = "service.exec_dir"
	spaExec          = "service.spa_exec"
)

// defaultConfig - config default value
var defaultConfig = Config{
	DataDir:           "./data",
	RootDir:           "./workspace",
	SpartaPythonTools: "/home/sparta/tools/pizza",
	ExecDir:           "/home/sparta/tools",
	SpaExec:           "/home/spa_",
}

// Config -
type Config struct {
	DataDir           string `toml:"data_dir" json:"data_dir"`                       // 数据存储目录
	RootDir           string `toml:"root_dir" json:"root_dir"`                       // 项目根目录
	SpartaPythonTools string `toml:"sparta_python_tools" json:"sparta_python_tools"` // python 脚本工具路径
	ExecDir           string `toml:"exec_dir" json:"exec_dir"`                       // 脚本执行目录
	SpaExec           string `toml:"spa_exec" json:"spa_exec"`                       // spa 目录
}

// SetDefaultConfig -
func SetDefaultConfig() {
	viper.SetDefault(dataDir, defaultConfig.DataDir)
	viper.SetDefault(rootDir, defaultConfig.RootDir)
	viper.SetDefault(spartaPytonTools, defaultConfig.SpartaPythonTools)
	viper.SetDefault(execDir, defaultConfig.ExecDir)
	viper.SetDefault(spaExec, defaultConfig.SpaExec)
}

// GetConfig -
func GetConfig() *Config {
	return &Config{
		DataDir:           viper.GetString(dataDir),
		RootDir:           viper.GetString(rootDir),
		SpartaPythonTools: viper.GetString(spartaPytonTools),
		ExecDir:           viper.GetString(execDir),
		SpaExec:           viper.GetString(spaExec),
	}
}
