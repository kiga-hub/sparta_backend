package service

import "github.com/spf13/viper"

const (
	datadir = "service.data_dir"
)

// 配置默认值 - 最低优先级
var defaultConfig = Config{
	Datadir: "./data",
}

// Config - 配置结构
type Config struct {
	Datadir string `toml:"datadir" json:"datadir"`
}

// SetDefaultConfig - 设置默认配置
func SetDefaultConfig() {
	viper.SetDefault(datadir, defaultConfig.Datadir)
}

// GetConfig - 获取当前配置
func GetConfig() *Config {
	return &Config{
		Datadir: viper.GetString(datadir),
	}
}
