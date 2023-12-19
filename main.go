package main

import (
	"context"
	"flag"
	"io"
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
	go hub.run()

	//
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/styles.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/css/styles.css")
	})
	router.GET("/login.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/css/login.css")
	})
	router.GET("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/login.html")
	})
	router.POST("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, _ := io.ReadAll(r.Body)
		log.Println(string(data))
		// TODO create auth middleware
		// email=?&password=?
		// TODO make call to authenticate and return token
		if func() bool { return true }() {
			http.ServeFile(w, r, "public/index.html")
		}
		// var msg json.RawMessage
		// if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		// 	log.Println("failed Decode:", err)
		// 	return
		// }
		// log.Printf("%#v\n", msg)
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
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("ListenAndServe:", err)
		}
	}()
	log.Println("started http server at", *addr)

	<-interruptCh
	log.Println("starting graceful shutdown")
}
