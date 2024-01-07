package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ChatService struct {
}

func (s *ChatService) StartNewChat(
	w http.ResponseWriter,
	r *http.Request,
	p httprouter.Params,
) {

}
