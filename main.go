package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	secureToken = ""

	startedMutex  = &sync.RWMutex{}
	serverStarted = false
)

func main() {

	// Get the security token from env var that will be used to verify calls from clients (X-API-TOKEN header)
	secureToken = os.Getenv("API_TOKEN")

	if secureToken == "" {
		log.Fatal("API_TOKEN env variable not set")
	}

	// This app is slow to start.... set global var to true when its finished its oh so important start up process
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

	// The endpoints for the app to listen on. /healthz and /readyz do not require authentication
	m.HandleFunc("/healthz", logHandler(healthzHandler()))
	m.HandleFunc("/readyz", logHandler(healthzHandler()))

	m.HandleFunc("/{uri:.*}", authHandler(logHandler(echoHandler())))

	log.Fatal(http.ListenAndServe(":8080", m))
}

// healthzHandler will return if the app is healthy and running
func healthzHandler() http.HandlerFunc {
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

// logHandler logs all request to the log / stdout
func logHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("request received. path=%v", request.URL)
		h(writer, request)
	}
}

// authHandler makes sure calls to specified handlers contain the correct security header, as configured by the env var API_TOKEN
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

// echoHandler will return a JSON object containing the requested URI and the FOO and BAR env vars from the server
// it will return even if the app is not marked as healthy in /healthz
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
	Path string            `json:"requestedPath"`
	Env  map[string]string `json:"env"`
}
