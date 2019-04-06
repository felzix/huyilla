package main

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func pingHandler(_ *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "pong"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
			if _, err := fmt.Fprint(w, "Signup successful!"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if err.Error()[:8] == "Player \"" {
			http.Error(w, err.Error(), http.StatusConflict)
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
			if _, err := fmt.Fprint(w, token); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		name, tokenId, _, err := readToken(engine.Secret, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		player, err := engine.World.Player(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if player == nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		} else {
			if _, currentTokenId, _, err := readToken(engine.Secret, player.Token); err == nil {
				if currentTokenId != tokenId {
					http.Error(w, "Old token", http.StatusForbidden)
					return
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := engine.LogOut(name); err == nil {
			if _, err := fmt.Fprint(w, "Logout successful!"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func userExistsHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exists, err := engine.UserExists(string(username)); err == nil {
			if _, err := fmt.Fprint(w, exists); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func playerHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if player, err := engine.World.Player(vars["name"]); err == nil {
			if blob, err := player.Marshal(); err == nil {
				if _, err := fmt.Fprint(w, blob); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (engine *Engine) Serve(errChan chan error) {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler(engine))

	r = mux.NewRouter()
	r.PathPrefix("/auth/")
	r.HandleFunc("signup", signupHandler(engine))
	r.HandleFunc("login", loginHandler(engine))
	r.HandleFunc("logout", logoutHandler(engine))
	r.HandleFunc("exists", userExistsHandler(engine))

	r = mux.NewRouter()
	r.PathPrefix("/world/")
	r.HandleFunc("player/{name}", playerHandler(engine))
	// http.HandleFunc("/chunk/", chunkHandler)
	// http.HandleFunc("/stats", statsHandler)

	errChan <- http.ListenAndServe(":8080", nil)
}
