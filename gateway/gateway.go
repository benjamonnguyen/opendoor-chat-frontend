package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

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
}
