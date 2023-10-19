package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// operation is a clean up function on shutting down
type ShutdownHook func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
// https://medium.com/tokopedia-engineering/gracefully-shutdown-your-go-application-9e7d5c73b5ac
func GracefulShutdown(ctx context.Context, timeout time.Duration, logger *logrus.Logger, ops map[string]ShutdownHook) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		println()
		logger.Info("Shutting down…")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			logger.Warnf("Timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				logger.Debugf("Cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					logger.Errorf("%s: Clean up failed: %s", innerKey, err.Error())
					return
				}

				logger.Debugf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
