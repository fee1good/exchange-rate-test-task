package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fee1good/exchange-rate-test-task/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := internal.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	go app.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-signalChan

	if err := app.Shutdown(); err != nil {
		panic(err)
	}
}
