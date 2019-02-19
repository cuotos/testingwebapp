package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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

	secureToken = ""
)

func main() {

	secureToken = os.Getenv("API_TOKEN")

	if secureToken == "" {
		log.Fatal("API_TOKEN env variable not set")
	}

	if os.Getenv("DEBUG") != "1" {
		log.Println("aligning the stars...")
		log.Println("starting service in...")

		for i := len(loadingMsgs) - 1; i >= 0; i-- {
			log.Printf("%v... %v", i+1, loadingMsgs[i])
			time.Sleep(time.Second * 2)
		}
	}

	m := mux.NewRouter()

	m.HandleFunc("/healthz", logHandler(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("ok"))
		return
	}))
	m.HandleFunc("/{uri}", authHandler(logHandler(echoHandler())))

	log.Println("service running")
	log.Fatal(http.ListenAndServe(":8080", m))
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
	Path string `json:"path"`
	Env map[string]string `json:"env"`
}