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
	"github.com/benjamonnguyen/opendoor-chat-frontend/chat"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
)

func main() {
	addr := flag.String("addr", ":3000", "http service address")
	flag.Parse()
	devlog.Enable(true)

	// graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	// config
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	// ws
	hub := chat.NewHub()
	go hub.Run()

	// server
	cl := &http.Client{
		Timeout: time.Minute,
	}
	srv := buildServer(cfg, *addr, hub, cl)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("ListenAndServe:", err)
		}
	}()
	log.Println("started http server at", *addr)

	<-interruptCh
	log.Println("starting graceful shutdown")
}
