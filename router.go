package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/benjamonnguyen/opendoor-chat-frontend/chat"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
	"github.com/benjamonnguyen/opendoor-chat-frontend/devlog"
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
	cl = &http.Client{
		Timeout: time.Minute,
	}
)

func buildServer(cfg config.Config, addr string, hub *chat.Hub) *http.Server {
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
	router.POST("/api/login", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, _ := io.ReadAll(r.Body)
		devlog.Println(string(data))
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
	router.POST("/api/signup", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// get params
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Str("route", "POST /api/signup").Msg("failed reading body")
			http.Error(w, "", 500)
			return
		}
		vals, err := url.ParseQuery(string(data))
		if err != nil {
			log.Error().Err(err).Str("route", "POST /api/signup").Msg("failed reading body")
			http.Error(w, "", 500)
			return
		}

		// create user
		data, err = json.Marshal(map[string]string{
			"firstName": vals.Get("first-name"),
			"lastName":  vals.Get("last-name"),
			"email":     vals.Get("email"),
			"password":  vals.Get("password"),
		})
		if err != nil {
			log.Error().
				Str("route", "POST /api/signup").
				Err(err).
				Msg("failed marshal")
			http.Error(w, "", 500)
			return
		}
		req, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodPost,
			cfg.BackendBaseUrl+"/user",
			bytes.NewReader(data),
		)
		if err != nil {
			log.Error().
				Str("route", "POST /api/signup").
				Err(err).
				Msg("failed constructing request")
			http.Error(w, "", 500)
			return
		}
		resp, err := cl.Do(req)
		if err != nil {
			log.Error().
				Str("route", "POST /api/signup").
				Err(err).
				Msg("failed request")
			http.Error(w, "", 500)
			return
		}

		// handle response
		var html string
		if resp.StatusCode == 201 {
			// TODO onboarding page
			html = `<dialog id="signup-modal" open><article>
			<header><Link aria-label="Close" class="close" hx-on:click="document.getElementById('signup-modal')?.remove();" /><strong>Thank You for Registering!</strong></header>
			<p>Welcome to Opendoor.chat!<br>Please verify your email.</p>
			</article></dialog>`
		} else {
			html = `<dialog id="signup-modal" open><article>
			<header><Link aria-label="Close" class="close" hx-on:click="document.getElementById('signup-modal')?.remove();" /></header>
			<p>Yikes! Look like something went wrong.</p>
			</article></dialog>`
		}
		w.Write([]byte(html))
		// TODO verification email
	})

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
	// check for auth bear token
	// make call to auth-svc
	// handle resp
}
