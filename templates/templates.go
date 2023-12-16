package templates

import (
	"html/template"
	"log"
)

var (
	ChatMessageTemplate *template.Template
)

func init() {
	t, err := template.New("chatMessage").
		Parse(`<div id="chat_messages" hx-swap-oob="beforeend">{{.}}</div>`)
	if err != nil {
		log.Fatalln("failed parsing chatMessage template:", err)
	}
	ChatMessageTemplate = t
}
