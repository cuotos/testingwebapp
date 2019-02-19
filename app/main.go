package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var (
	loadingMsgs = []string{
		"fixing the world",
		"calculating bath bomb ballistic trajectory",
		"giving wally time to hide",
		"milking the almond herd",
		"mixing genetic pool",
	}

	debug = false
)

func main() {

	if !debug {
		log.Println("aligning the stars...")
		log.Println("starting service in...")

		for i := len(loadingMsgs) - 1; i >= 0; i-- {
			log.Printf("%v... %v", i+1, loadingMsgs[i])
			time.Sleep(time.Second * 2)
		}
	}

	m := mux.NewRouter()

	m.HandleFunc("/{uri}", echoHandler())

	log.Println("service running")
	log.Fatal(http.ListenAndServe(":8080", m))
}

func echoHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		requestUri, _ := vars["uri"]

		writer.Write([]byte(requestUri))
		return
	}
}