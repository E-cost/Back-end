package handlers

import "github.com/julienschmidt/httprouter"

type ItemHandler interface {
	SaveItem(router *httprouter.Router)
}
