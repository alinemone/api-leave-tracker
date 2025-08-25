package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"leave/bootstrap"

	"github.com/getsentry/sentry-go"
	"github.com/goravel/framework/facades"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start http server by facades.Route().
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route Run error: %v", err)
		}
	}()

	// Listen for the OS signal
	go func() {
		<-quit
		if err := facades.Route().Shutdown(); err != nil {
			facades.Log().Errorf("Route Shutdown error: %v", err)
		}

		os.Exit(0)
	}()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              facades.Config().GetString("SENTRY_DSN", "https://692a097d15fca0a3122d7a7b2f0ed492@sentry.basalam.com/504"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	select {}
}
