package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	secureToken = ""

	startedMutex  = &sync.RWMutex{}
	serverStarted = false
)

func main() {

	secureToken = os.Getenv("API_TOKEN")

	if secureToken == "" {
		log.Fatal("API_TOKEN env variable not set")
	}

	go func() {
		for i := 5; i > 0; i-- {
			log.Printf("starting in %v...", i)
			time.Sleep(time.Second * 5)
		}

		startedMutex.Lock()
		serverStarted = true
		startedMutex.Unlock()
		log.Println("service running")
	}()

	m := mux.NewRouter()

	m.HandleFunc("/healthz", logHandler(healthHandler()))
	m.HandleFunc("/readyz", logHandler(healthHandler()))
	m.HandleFunc("/{uri}", authHandler(logHandler(echoHandler())))

	log.Fatal(http.ListenAndServe(":8080", m))
}

func healthHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		startedMutex.RLock()
		started := serverStarted
		startedMutex.RUnlock()

		if !started {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("error: app not ready yet. normally takes approx 20s"))
			return
		} else {
			writer.Write([]byte("ok"))
			return
		}
	}
}

func logHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("request received. path=%v", request.URL)
		h(writer, request)
	}
}

func authHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		suppliedAuthToken := request.Header.Get("X-API-TOKEN")

		if suppliedAuthToken == secureToken {
			h(writer, request)
		} else {
			log.Println("unauthorized request")
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("error: unauthorized"))
			return
		}
	}
}

func echoHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		requestUri, _ := vars["uri"]

		envVars := make(map[string]string)
		envVars["FOO"] = os.Getenv("FOO")
		envVars["BAR"] = os.Getenv("BAR")

		res := responseBody{requestUri, envVars}

		resJson, _ := json.Marshal(res)

		writer.Write([]byte(resJson))
		return
	}
}

type responseBody struct {
	Path string            `json:"path"`
	Env  map[string]string `json:"env"`
}
