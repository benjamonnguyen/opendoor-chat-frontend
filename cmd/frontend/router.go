package main

import (
	"net/http"

	app "github.com/benjamonnguyen/opendoor-chat-frontend"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
	"github.com/benjamonnguyen/opendoor-chat-frontend/ws"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func buildServer(
	cfg config.Config,
	addr string,
	hub *ws.Hub,
	cl *http.Client,
	chatSvc *app.ChatService,
) *http.Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router := httprouter.New()
	//
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Hello, World!"))
	})
	// App pages
	router.GET("/app", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if !isAuthorized(r) {
			http.Redirect(w, r, "/app/login", http.StatusFound)
			return
		}
		http.ServeFile(w, r, "public/app.html")
	})
	router.GET("/app/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/login.html")
	})
	router.GET("/app/signup", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		http.ServeFile(w, r, "public/signup.html")
	})
	// TODO /app/demo get demo data to populate UI and allow user to click around, but don't allow mutation

	// API Gateway interfaces with backend
	gateway := app.NewApiGateway(cl, cfg)
	router.POST("/api/login", gateway.LogIn)
	router.POST("/api/signup", gateway.SignUp)

	// Chat
	router.GET("/chat/start-new-chat", chatSvc.StartNewChat)

	// WS
	router.GET("/ws", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("failed ws upgrade")
			return
		}
		hub.Register(ws.NewClient(hub, conn))
	})

	// CSS
	router.GET(
		"/css/:file",
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			http.ServeFile(w, r, "public/css/"+p.ByName("file"))
		},
	)

	// Assets
	router.GET(
		"/assets/*filepath",
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			http.ServeFile(w, r, "public/"+p.ByName("filepath"))
		},
	)

	//
	n := negroni.Classic()
	n.UseHandler(router)

	//
	return &http.Server{
		Addr:    addr,
		Handler: n,
	}
}

func isAuthorized(r *http.Request) bool {
	c, err := r.Cookie("OPENDOOR_CHAT_TOKEN")
	if err != nil {
		return false
	}
	// TODO isAuthorized impl
	return c.Value != ""
}
