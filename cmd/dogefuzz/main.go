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
	"github.com/dogefuzz/dogefuzz/job"
	"github.com/dogefuzz/dogefuzz/listener"
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

	// Run job scheduler
	scheduler := job.NewJobScheduler(cfg)
	scheduler.Start()

	// Run listener manager
	listenerManager := listener.NewManager(cfg)
	listenerManager.Start()

	waitForInterrupt(server, scheduler, listenerManager)
}

func waitForInterrupt(svr api.Server, scheduler job.Scheduler, manager listener.Manager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	manager.Shutdown()
	log.Println("shutting down listener manager")

	ctx := scheduler.Shutdown()
	<-ctx.Done()
	log.Println("shutting down job scheduler")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Fatal("server Shutdown error", err)
	}
	<-ctx.Done()
	log.Println("shutting down server")

	os.Exit(0)
}
