package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

var (
	addr      = flag.String("addr", ":3000", "http service address")
	devLogger = NewDevLogger(true)
)

func main() {
	flag.Parse()

	// graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	//
	hub := newHub()
	defer func() {
		for cl := range hub.clients {
			cl.conn.Close()
		}
	}()
	go hub.run()

	//
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/styles.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/styles.css")
	})
	router.GET("/ws", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		serveWs(hub, w, r)
	})

	//
	n := negroni.Classic()
	n.UseHandler(router)
	srv := &http.Server{
		Addr:    *addr,
		Handler: n,
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("ListenAndServe:", err)
		}
	}()

	<-interruptCh
}
