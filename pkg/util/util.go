package util

import (
	"context"
	"github.com/go-logr/logr"
	"os"
)

// ContextFromStopChannel instantiates a context that is open as long as the stopCh is open.
// It will return context.ErrCanceled as soon as the stopCh is closed.
func ContextFromStopChannel(stopCh <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		<-stopCh
	}()
	return ctx
}

// LogErrorAndExit logs the error and exits with code 1.
func LogErrorAndExit(log logr.Logger, err error, msg string, keysAndValues ...interface{}) {
	log.Error(err, msg, keysAndValues)
	os.Exit(1)
}
