package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

func (conf *Conf) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(conf.Home))
	mux.Get("/ws", http.HandlerFunc(WsEndpoint))
	return mux
}
