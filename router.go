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
	// Pages
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/login.html")
	})
	router.GET("/signup", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/signup.html")
	})
	// CSS
	router.GET("/styles.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/css/styles.css")
	})
	router.GET("/login.css", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/css/login.css")
	})
	//
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
