package main

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"net/http"
)

func pingHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func signupHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blob []byte
		if _, err := r.Body.Read(blob); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var auth types.Auth
		if err := auth.Unmarshal(blob); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := engine.SignUp(auth.Name, string(auth.Password)); err == nil {
			fmt.Fprint(w, "Signup successful!")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (engine *Engine) Serve(errChan chan error) {
	http.HandleFunc("/auth/signup", signupHandler(engine))
	http.HandleFunc("/ping", pingHandler(engine))
	// http.HandleFunc("/auth/login", loginHandler)
	// http.HandleFunc("/player/", playerHandler)
	// http.HandleFunc("/chunk/", chunkHandler)
	// http.HandleFunc("/stats", statsHandler)

	errChan <- http.ListenAndServe(":8080", nil)
}
