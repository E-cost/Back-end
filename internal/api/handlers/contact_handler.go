package handlers

import "github.com/julienschmidt/httprouter"

type ContactHandler interface {
	SaveContact(router *httprouter.Router)
}
