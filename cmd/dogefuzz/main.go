package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogefuzz/dogefuzz/api"
	"github.com/dogefuzz/dogefuzz/config"
)

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Couldn't load config")
		panic(err)
	}

	// Run server
	server := api.NewServer(cfg)
	if err = server.Start(); err != nil {
		log.Fatal("Couldn't start server")
		panic(err)
	}
	waitForInterrupt(server)
}

func waitForInterrupt(svr api.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown error", err)
	}

	<-ctx.Done()
	log.Println("Shutting down server")
	os.Exit(0)
}
