package main

import (
	"net/http"

	"github.com/benjamonnguyen/opendoor-chat-frontend/chat"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
	"github.com/benjamonnguyen/opendoor-chat-frontend/gateway"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func buildServer(cfg config.Config, addr string, hub *chat.Hub, cl *http.Client) *http.Server {
	// App pages
	router := httprouter.New()
	router.GET("/app", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO authenticate
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
	gateway := gateway.NewApiGateway(cl, cfg)
	router.POST("/api/login", gateway.LogIn)
	router.POST("/api/signup", gateway.SignUp)

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
		log.Error().Err(err).Msg("failed ws upgrade")
		return
	}
	hub.Register(chat.NewClient(hub, conn))
}

func auth() {
	// TODO check for auth bear token
	// make call to auth-svc
	// handle resp
}
