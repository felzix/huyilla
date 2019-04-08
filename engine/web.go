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
		if r.Body == nil {
			http.Error(w, "Must supply body", http.StatusBadRequest)
			return
		}

		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var auth types.Auth
		if err := auth.Unmarshal(blob); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		if r.Body == nil {
			http.Error(w, "Must supply body", http.StatusBadRequest)
			return
		}

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
		name := vars["name"]

		if name == "" {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		player, err := engine.World.Player(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if player == nil {
			http.Error(w, fmt.Sprintf(`No such player "%s"`, name), http.StatusNotFound)
			return
		}

		sentPlayer := types.Player{
			EntityId: player.EntityId,
			Name: player.Name,
		}

		if blob, err := sentPlayer.Marshal(); err == nil {
			w.Header().Set("Content-Type", "application/protobuf")
			if _, err := w.Write(blob); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Router(engine *Engine) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler(engine)).Methods("GET")

	r.HandleFunc("/auth/signup", signupHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/login", loginHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/logout", logoutHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/exists", userExistsHandler(engine)).Methods("GET")

	r.HandleFunc("/world/player/{name}", playerHandler(engine)).Methods("GET")
	// http.HandleFunc("/chunk/", chunkHandler)
	// http.HandleFunc("/stats", statsHandler)

	return r
}


func (engine *Engine) Serve(errChan chan error) {
	Router(engine)
	errChan <- http.ListenAndServe(":8080", nil)
}
