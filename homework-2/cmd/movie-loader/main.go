package main

import (
	"context"
	"homework-2/internal/cli"
	"os/signal"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	rootCmd := cli.RootCmd(ctx)

	if err := rootCmd.Execute(); err != nil {
		return
	}
}
