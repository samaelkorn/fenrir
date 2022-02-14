package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	internalhttp "github.com/samaelkorn/fenrir/internal/http"
)

func main() {
	server := internalhttp.NewServer("localhost", "6078")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Print("failed to stop http server: " + err.Error())
		}
	}()

	log.Print("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Print("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

}
