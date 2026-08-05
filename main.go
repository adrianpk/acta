// Command acta runs the local chronicle client and service.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"forge.adrianpk.com/adrianpk/acta/internal/acta"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := acta.DefaultConfig()
	if err == nil {
		err = acta.NewApp(config, os.Stdout, os.Stderr).Run(ctx, os.Args[1:])
	}

	if err != nil {
		log.New(os.Stderr, "", 0).Println(err)
		os.Exit(1)
	}
}
