package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/benjamonnguyen/gootils/devlog"
	app "github.com/benjamonnguyen/opendoor-chat-frontend"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
	"github.com/benjamonnguyen/opendoor-chat-frontend/ws"
)

func main() {
	addr := flag.String("addr", ":3000", "http service address")
	flag.Parse()
	devlog.Init(true, nil)

	// graceful shutdown setup
	ctx, cancel := context.WithCancel(context.Background())
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		<-interruptCh
		cancel()
	}()

	// config
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	// services
	chatSvc := new(app.ChatService)

	// ws
	hub := ws.NewHub()
	go hub.Run(ctx)

	// server
	cl := &http.Client{
		Timeout: time.Minute,
	}
	srv := buildServer(cfg, *addr, hub, cl, chatSvc)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("ListenAndServe:", err)
		}
	}()
	log.Println("started http server at", *addr)

	<-ctx.Done()
	// graceful shutdown
	start := time.Now()
	shtudownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	srv.Shutdown(shtudownCtx)

	log.Printf("completed graceful shutdown after %s", time.Since(start))
}
