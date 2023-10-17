package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"mqtt_web_api/api"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SPAHandler struct {
	StaticPath string
	IndexPath  string
}

func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.StaticPath, r.URL.Path)

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.Now()

		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api") {
			w.Header().Set("Content-Type", "application/json")
			origin := r.Header.Get("Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		handler.ServeHTTP(w, r)
		fmt.Printf("[%s] (%s) %s in %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(timer))
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	api.MakeHandler(s)
	r.Use(Middleware)

	spa := SPAHandler{StaticPath: "dist", IndexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	fmt.Println("Listening on port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
