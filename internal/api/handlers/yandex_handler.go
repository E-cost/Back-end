package handlers

import "github.com/julienschmidt/httprouter"

type YandexHandler interface {
	SaveYandexStorage(router *httprouter.Router)
}
