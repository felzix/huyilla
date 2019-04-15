package engine

import (
	"fmt"
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
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
			if _, err := w.Write([]byte(token)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, tokenId, _, err := engine.authenticate(w, r)
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
		vars := mux.Vars(r)
		username := vars["name"]

		if exists, err := engine.UserExists(username); err == nil {
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

		entity, err := engine.World.Entity(player.EntityId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if entity == nil {
			http.Error(
				w,
				fmt.Sprintf("Player's entity %d does not exist", player.EntityId),
				http.StatusInternalServerError)
			return
		}

		thisUser, _, _, err := engine.authenticate(w, r)
		if err != nil && err.Error() == "must specify token" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var toSend *types.Entity
		if thisUser == player.Name {
			toSend = entity
		} else {
			toSend = &types.Entity{
				Id: entity.Id,
			}
		}

		if blob, err := toSend.Marshal(); err == nil {
			w.Header().Set("Content-Type", "application/protobuf")
			if _, err := w.Write(blob); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func chunkHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, _, _, err := engine.authenticate(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		chunkPoint, err := stringToPoint(vars["x"], vars["y"], vars["z"])
		if err != nil {
			http.Error(w, "bad url; must be /world/chunk/{x}/{y}/{z}", http.StatusBadRequest)
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

		playerEntity, err := engine.World.Entity(player.EntityId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if gridDistance(playerEntity.Location.Chunk, chunkPoint) > constants.ACTIVE_CHUNK_RADIUS {
			http.Error(w, "can only load nearby chunks", http.StatusForbidden)
			return
		}

		chunk, err := engine.World.Chunk(chunkPoint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if chunk == nil {
			http.Error(w, fmt.Sprintf(`No such Chunk "%s"`, pointToString(chunkPoint)), http.StatusNotFound)
			return
		}

		entities := make([]*types.Entity, len(chunk.Entities))
		for i := 0; i < len(chunk.Entities); i++ {
			id := chunk.Entities[i]
			entity, err := engine.World.Entity(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			entities[i] = entity
		}

		toSend := &types.ChunkDetail{
			Voxels: chunk.Voxels,
			Compound: chunk.Compound,
			Entities: entities,
			Items: chunk.Items,
		}

		if blob, err := toSend.Marshal(); err == nil {
			w.Header().Set("Content-Type", "application/protobuf")
			if _, err := w.Write(blob); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func worldAgeHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		age, err := engine.World.Age()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ticks := []byte(strconv.FormatUint(age.Ticks, 10))

		if _, err := w.Write(ticks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func actHandler(engine *Engine) http.HandlerFunc {
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

		var action types.Action
		if err := action.Unmarshal(blob); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		engine.RegisterAction(&action)
	}
}


func Router(engine *Engine) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler(engine)).Methods("GET")

	r.HandleFunc("/auth/signup", signupHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/login", loginHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/logout", logoutHandler(engine)).Methods("POST")
	r.HandleFunc("/auth/exists/{name}", userExistsHandler(engine)).Methods("GET")

	r.HandleFunc("/world/age", worldAgeHandler(engine)).Methods("GET")
	r.HandleFunc("/world/player/{name}", playerHandler(engine)).Methods("GET")
	r.HandleFunc("/world/chunk/{x}/{y}/{z}", chunkHandler(engine)).Methods("GET")
	r.HandleFunc("/world/act", actHandler(engine)).Methods("POST")
	// http.HandleFunc("/stats", statsHandler)

	return r
}

func (engine *Engine) Serve(port int, errChan chan error) *http.Server {
	r := Router(engine)
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		// returns ErrServerClosed on graceful close
		errChan <- srv.ListenAndServe()
	}()

	return srv
}

func (engine *Engine) authenticate(w http.ResponseWriter, r *http.Request) (string, string, int64, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", "", 0, errors.New("must specify token")
	}
	if name, tokenId, expiry, err := readToken(engine.Secret, token); err == nil {
		return name, tokenId, expiry, nil
	} else {
		return "", "", 0, err
	}
}
