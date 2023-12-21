package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/benjamonnguyen/gootils/httputil"
	"github.com/benjamonnguyen/opendoor-chat-frontend/config"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type ApiGateway struct {
	cl  *http.Client
	cfg config.Config
}

func NewApiGateway(cl *http.Client, cfg config.Config) *ApiGateway {
	return &ApiGateway{
		cl:  cl,
		cfg: cfg,
	}
}

func (a *ApiGateway) LogIn(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get params
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("route", "POST /api/login").Msg("failed reading body")
		http.Error(w, "", 500)
		return
	}
	vals, err := url.ParseQuery(string(data))
	if err != nil {
		log.Error().Err(err).Str("route", "POST /api/login").Msg("failed parsing query")
		http.Error(w, "", 500)
		return
	}

	// authenticate
	data, err = json.Marshal(map[string]string{
		"email":    vals.Get("email"),
		"password": vals.Get("password"),
	})
	if err != nil {
		log.Error().
			Str("route", "POST /api/login").
			Err(err).
			Msg("failed marshal")
		http.Error(w, "", 500)
		return
	}
	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		a.cfg.BackendBaseUrl+"/user/authenticate",
		bytes.NewReader(data),
	)
	if err != nil {
		log.Error().
			Str("route", "POST /api/login").
			Err(err).
			Msg("failed constructing request")
		http.Error(w, "", 500)
		return
	}
	resp, err := a.cl.Do(req)
	if err != nil {
		log.Error().
			Str("route", "POST /api/login").
			Err(err).
			Msg("failed request")
		http.Error(w, "", 500)
		return
	}

	// handle response
	if resp.StatusCode == 200 {
		token, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().
				Str("route", "POST /api/login").
				Err(err).
				Msg("failed reading repsonse body")
			http.Error(w, "", 500)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "Authorization",
			Value:   string(token),
			Expires: time.Now().Add(24 * time.Hour * 60),
		})
		w.Header().Set("HX-Redirect", "/app")
	} else if httputil.Is4xx(resp.StatusCode) {
		w.Write([]byte(`<div id="login-status"><small id="login-status-text" style="color: #FF6161;">
		The email and/or password you entered are not correct.</small></div>`))
	} else {
		w.Write([]byte(`<div id="login-status"> <small id="login-status-text" style="color: #FF6161;">
		Something went wrong. Please wait a moment and try again.</small></div>`))
	}
}

func (a *ApiGateway) SignUp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get params
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Str("route", "POST /api/signup").Msg("failed reading body")
		http.Error(w, "", 500)
		return
	}
	vals, err := url.ParseQuery(string(data))
	if err != nil {
		log.Error().Err(err).Str("route", "POST /api/signup").Msg("failed parsing query")
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
		a.cfg.BackendBaseUrl+"/user",
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
	resp, err := a.cl.Do(req)
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
		html = `<div id="login-status"><small id="login-status-text">
		You're registered! Please verify your email.</small></div>`
	} else if resp.StatusCode == http.StatusConflict {
		html = `<div id="login-status"><small id="login-status-text" style="color: #FF6161;">
		This email is already in use. <a href="/app/login>Log in?</a></small></div>`
	} else {
		log.Error().Str("route", "POST /api/signup").Msg(resp.Status)
		html = `<div id="login-status"><small id="login-status-text" style="color: #FF6161;">
		Something went wrong. Please wait a moment and try again.</small></div>`
	}
	w.Write([]byte(html))
	// TODO verification email
}
