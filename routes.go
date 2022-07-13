package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

func (conf *Conf) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(conf.Home))
	mux.Get("/ws", http.HandlerFunc(WsEndpoint))
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
