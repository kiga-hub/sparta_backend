package service

import "github.com/spf13/viper"

const (
	dataDir = "service.data_dir"
)

// defaultConfig - config default value
var defaultConfig = Config{
	DataDir: "./data",
}

// Config -
type Config struct {
	DataDir string `toml:"datadir" json:"datadir"`
}

// SetDefaultConfig -
func SetDefaultConfig() {
	viper.SetDefault(dataDir, defaultConfig.DataDir)
}

// GetConfig -
func GetConfig() *Config {
	return &Config{
		DataDir: viper.GetString(dataDir),
	}
}
