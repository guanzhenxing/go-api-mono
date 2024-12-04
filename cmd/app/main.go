package main

import (
	"fmt"
	"os"

	"go-api-mono/internal/app"
)

func main() {
	// 获取配置文件路径
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "configs/config.yaml"
	}

	// 创建并运行应用程序
	opts := app.Options{
		ConfigFile: configFile,
		LogLevel:   os.Getenv("LOG_LEVEL"),
		DevMode:    os.Getenv("DEV_MODE") == "true",
	}

	if err := app.Run(opts); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}
