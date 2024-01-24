package shutdown

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ShutdownGroup interface {
	Register(string, func(context.Context) error)
	Wait(ctx context.Context) <-chan struct{}
	WaitBlocking(ctx context.Context)
}

type shutdownGroup struct {
	logger  *slog.Logger
	timeout time.Duration
	signals []os.Signal
	ops     map[string]func(context.Context) error
}

func NewShutdownGroup(logger *slog.Logger, timeout time.Duration, signals ...os.Signal) ShutdownGroup {
	return &shutdownGroup{logger: logger, timeout: timeout, ops: map[string]func(context.Context) error{}, signals: signals}
}

func (sg *shutdownGroup) Register(name string, operation func(ctx context.Context) error) {
	sg.ops[name] = operation
}

func (sg *shutdownGroup) WaitBlocking(ctx context.Context) {
	wait := sg.Wait(ctx)
	<-wait
}

func (sg *shutdownGroup) Wait(ctx context.Context) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, sg.signals...)
		sig := <-s
		sg.logger.Info("shutdown signal received", "signal", sig.String(), "timeout", sg.timeout.String())
		timeoutFunc := time.AfterFunc(sg.timeout, func() {
			sg.logger.Error("shutdown timeout reached, forced exit", "timeout", sg.timeout.String())
			os.Exit(0)
		})
		defer timeoutFunc.Stop()
		var wg sync.WaitGroup
		for key, op := range sg.ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()
				sg.logger.Info("shutting down", "operation", innerKey)
				if err := innerOp(ctx); err != nil {
					sg.logger.Error("shutdown failed", "operation", innerKey, "error", err.Error())
					return
				}
				sg.logger.Info("graceful shutdown", "operation", innerKey)
			}()
		}
		wg.Wait()
		close(wait)
	}()
	return wait
}
