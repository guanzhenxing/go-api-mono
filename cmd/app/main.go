package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"go-api-mono/internal/app"
	"go-api-mono/internal/pkg/config"
)

// main 是应用程序的入口点
// 它负责：
// 1. 加载配置
// 2. 初始化应用
// 3. 启动服务
// 4. 处理优雅关闭
func main() {
	// 加载配置文件，如果加载失败则直接退出
	app, err := app.New(config.MustLoad())
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个可以监听系统信号的上下文
	// 当收到 SIGINT 或 SIGTERM 信号时，会触发上下文取消
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 在单独的 goroutine 中启动应用
	// 这样可以不阻塞主 goroutine，使其能够处理关闭信号
	go func() {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// 等待上下文取消（即收到关闭信号）
	<-ctx.Done()

	// 开始优雅关闭流程
	// 使用新的上下文来避免已取消的上下文影响关闭流程
	if err := app.Stop(context.Background()); err != nil {
		log.Printf("Failed to stop application: %v", err)
	}
}
