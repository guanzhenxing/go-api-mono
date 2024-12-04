package configs

import (
	"embed"
	"fmt"
	"os"
)

//go:embed *.yaml
var configFS embed.FS

// GetConfigFile 从嵌入的文件系统中读取配置文件
func GetConfigFile(env string) ([]byte, error) {
	// 如果环境变量中指定了配置文件路径，优先使用外部配置
	if configPath := os.Getenv("CONFIG_FILE"); configPath != "" {
		return os.ReadFile(configPath)
	}

	// 根据环境确定配置文件名
	var filename string
	if env == "" {
		filename = "config.yaml"
	} else {
		filename = fmt.Sprintf("config.%s.yaml", env)
	}

	// 从嵌入的文件系统中读取
	return configFS.ReadFile(filename)
}
