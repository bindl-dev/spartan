package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"

	"github.com/bindl-dev/sparta"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err.Error())
		os.Exit(1)
	}
}

//go:embed sparta.json
var conf []byte

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	s := sparta.NewServer(ctx, conf)

	go func() {
		<-ctx.Done()
		_ = s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}
