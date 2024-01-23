package service

import "github.com/spf13/viper"

const (
	dataDir          = "service.data_dir"
	rootDir          = "service.root_dir"
	spartaPytonTools = "service.sparta_python_tools"
	scriptDir        = "service.script_dir"
	spaExec          = "service.spa_exec"
)

// defaultConfig - config default value
var defaultConfig = Config{
	DataDir:           "./data",
	RootDir:           "./workspace",
	SpartaPythonTools: "/home/sparta/tools/pizza",
	ScriptDir:         "/home/sparta/tools",
	SpaExec:           "/home/spa_",
}

// Config -
type Config struct {
	DataDir           string `toml:"data_dir" json:"data_dir"`                       // data dir
	RootDir           string `toml:"root_dir" json:"root_dir"`                       // workspace dir
	SpartaPythonTools string `toml:"sparta_python_tools" json:"sparta_python_tools"` // python package dir
	ScriptDir         string `toml:"exec_dir" json:"exec_dir"`                       // script dir
	SpaExec           string `toml:"spa_exec" json:"spa_exec"`                       // spa exec
}

// SetDefaultConfig -
func SetDefaultConfig() {
	viper.SetDefault(dataDir, defaultConfig.DataDir)
	viper.SetDefault(rootDir, defaultConfig.RootDir)
	viper.SetDefault(spartaPytonTools, defaultConfig.SpartaPythonTools)
	viper.SetDefault(scriptDir, defaultConfig.ScriptDir)
	viper.SetDefault(spaExec, defaultConfig.SpaExec)
}

// GetConfig -
func GetConfig() *Config {
	return &Config{
		DataDir:           viper.GetString(dataDir),
		RootDir:           viper.GetString(rootDir),
		SpartaPythonTools: viper.GetString(spartaPytonTools),
		ScriptDir:         viper.GetString(scriptDir),
		SpaExec:           viper.GetString(spaExec),
	}
}
