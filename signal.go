package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"
)

func Shutdown(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		Logger.Log("get", "signal", "signal", s)
		cancelFunc()
	}()
}
