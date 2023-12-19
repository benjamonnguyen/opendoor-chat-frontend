package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/benjamonnguyen/opendoor-chat-frontend/chat"
	"github.com/benjamonnguyen/opendoor-chat-frontend/devlog"
)

var (
	addr = flag.String("addr", ":3000", "http service address")
)

func main() {
	flag.Parse()
	devlog.Enable(true)

	// graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	//
	hub := chat.NewHub()
	go hub.Run()

	//
	srv := buildServer(*addr, hub)
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
