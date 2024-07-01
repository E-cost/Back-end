package handlers

import "github.com/julienschmidt/httprouter"

type OrderHandler interface {
	SaveOrder(router *httprouter.Router)
}
