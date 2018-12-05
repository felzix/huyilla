package main

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"io/ioutil"
	"net/http"
	"time"
)

func pingHandler(_ *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func signupHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
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

func loginHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var auth types.Auth
		if err := auth.Unmarshal(blob); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if token, err := engine.LogIn(auth.Name, string(auth.Password)); err == nil {
			fmt.Fprint(w, token)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		name, exp, err := readToken(engine.Secret, token)
		fmt.Println(token)
		fmt.Println(name)
		fmt.Println(exp)
		fmt.Println(err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exp > time.Now().Unix() {
			// TODO figure out auto-logout so the token isn't needed... would be bad to have a user stuck
			//      logged-in but unable to do anything, even log out
			//      ALSO note that this line shouldn't ever be hit because readToken calls Valid which throws on this
			//      condition
		}

		if err := engine.LogOut(name); err == nil {
			fmt.Fprint(w, "Logout successful!")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (engine *Engine) Serve(errChan chan error) {
	http.HandleFunc("/ping", pingHandler(engine))

	http.HandleFunc("/auth/signup", signupHandler(engine))
	http.HandleFunc("/auth/login", loginHandler(engine))
	http.HandleFunc("/auth/logout", logoutHandler(engine))

	// http.HandleFunc("/auth/login", loginHandler)
	// http.HandleFunc("/player/", playerHandler)
	// http.HandleFunc("/chunk/", chunkHandler)
	// http.HandleFunc("/stats", statsHandler)

	errChan <- http.ListenAndServe(":8080", nil)
}
