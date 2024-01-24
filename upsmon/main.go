package main

import (
	"context"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/enolgor/winpower/upsmon/service"
	"github.com/enolgor/winpower/upsmon/shutdown"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: service.LogLevel,
	}))
	upsmon := service.NewUpsMon(logger)
	go upsmon.Start()
	shutdown := shutdown.NewShutdownGroup(logger, 5*time.Second, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	shutdown.Register("UpsMon", upsmon.Shutdown)
	shutdown.WaitBlocking(context.Background())
}
