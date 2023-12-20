package main

import (
	"io"
	"log"
	"net/http"

	"github.com/benjamonnguyen/opendoor-chat-frontend/chat"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func buildServer(addr string, hub *chat.Hub) *http.Server {
	// App pages
	router := httprouter.New()
	router.GET("/app", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/app/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/login.html")
	})
	router.GET("/app/signup", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/signup.html")
	})
	// TODO /app/demo get demo data to populate UI and allow user to click around, but don't allow mutation

	// CSS
	router.GET(
		"/css/styles.css",
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			http.ServeFile(w, r, "public/css/styles.css")
		},
	)
	router.GET("/css/login.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/css/login.css")
	})

	// API
	router.POST("/api/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, _ := io.ReadAll(r.Body)
		log.Println(string(data))
		// TODO create auth middleware
		// email=?&password=?
		// TODO make call to authenticate and return token
		// salt and hash (sha256)
		// validate length and strength for signup
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
	// router.POST("/api/sign")

	// WS
	router.GET("/ws", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		serveWs(hub, w, r)
	})

	//
	n := negroni.Classic()
	n.UseHandler(router)

	//
	return &http.Server{
		Addr:    addr,
		Handler: n,
	}
}

func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.Register(chat.NewClient(hub, conn))
}

func auth() {
	// check for auth bear token
	// make call to auth-svc
	// handle resp
}
