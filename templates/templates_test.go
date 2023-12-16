package templates_test

import (
	"bytes"
	"testing"

	"github.com/benjamonnguyen/opendoor-chat-frontend/templates"
)

func TestChatMessageTemplate(t *testing.T) {
	buf := new(bytes.Buffer)
	templates.ChatMessageTemplate.Execute(buf, "test")
	want := `<div id="chat_messages" hx-swap-oob="beforeend">test</div>`
	got := buf.String()
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
