package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"go-api-mono/internal/app"
	"go-api-mono/internal/pkg/config"
)

func main() {
	app, err := app.New(config.MustLoad())
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	if err := app.Stop(context.Background()); err != nil {
		log.Printf("Failed to stop application: %v", err)
	}
}
