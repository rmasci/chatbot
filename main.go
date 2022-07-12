package main

import (
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
)

// Todo come up with a JSON, Yaml, or TOML config file.
type Conf struct {
	Port    string
	Rivedir string
	Host    string
	Botname string
}

func main() {
	var conf Conf
	pflag.StringVarP(&conf.Port, "port", "p", "8080", "Port to listen on")
	pflag.StringVarP(&conf.Rivedir, "rivescript", "r", "rive", "Rivescript directory.")
	// Todo: get this working with https -- this is just a test.
	pflag.StringVarP(&conf.Host, "host", "h", "10.1.1.130", "10.1.1.130 Typically your ip address.")
	pflag.StringVarP(&conf.Botname, "bot", "b", "chatbot", "Chatbot name")
	pflag.Parse()
	mux := conf.routes()
	log.Println("Starting channel listener")
	go ListenToWsChannel(conf)
	log.Println("Starting Web Server On Port", conf.Port)

	err := http.ListenAndServe(":"+conf.Port, mux)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
