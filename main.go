package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vox-imperialis/handlers"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	Load()

	cfg := Get()
	log.Printf("Vox Imperialis starting — JID: %s", cfg.JID)

	d := NewDispatcher()
	d.Register("help", handlers.Help)
	d.Register("status", handlers.Status)
	d.Register("sensors", handlers.Sensors)
	d.Register("service", handlers.NewServiceHandler(cfg.AllowedServices))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewXMPPClient(cfg, d)

	go func() {
		if err := client.Connect(ctx); err != nil && ctx.Err() == nil {
			log.Printf("xmpp: client terminated unexpectedly: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("received signal %s, shutting down", sig)
		cancel()
	case <-ctx.Done():
	}

	log.Println("Vox Imperialis stopped")
}
