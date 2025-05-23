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
		Log.Debug("web:ping")
		if _, err := fmt.Fprintf(w, "pong"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func signupHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("web:signup")
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
		Log.Debug("web:login")
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
		Log.Debug("web:logout")
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
		Log.Debugf("web:userExists: username=%v", username)

		if exists, err := engine.UserExists(username); err == nil {
			Log.Debugf("web: userExists: username=%v exists=%v", username, exists)
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

		Log.Debugf("web:player: name=%v", name)

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
		Log.Debug("web:chunk")
		name, _, _, err := engine.authenticate(w, r)
		if err != nil {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		center, err := types.StringToPoint(vars["x"], vars["y"], vars["z"])
		if err != nil {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, "bad url: must be /world/chunk/{x}/{y}/{z}", http.StatusBadRequest)
			return
		}

		var radius uint64
		radii, ok := r.URL.Query()["radius"]
		if ok && len(radii) == 1 {
			radius, err = strconv.ParseUint(radii[0], 10, 64)
			if err != nil {
				Log.Errorf("web:chunk: error=%v", err)
				http.Error(w, "bad url param: radius must be a positive integer", http.StatusBadRequest)
				return
			}
		} else {
			radius = 0
		}

		player, err := engine.World.Player(name)
		if err != nil {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if player == nil {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, fmt.Sprintf(`No such player "%s"`, name), http.StatusNotFound)
			return
		}

		playerEntity, err := engine.World.Entity(player.EntityId)
		if err != nil {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mostDistanceChunkPoint := center.Clone()
		mostDistanceChunkPoint.X += int64(radius)
		if playerEntity.Location.Chunk.GridDistance(mostDistanceChunkPoint) > constants.ACTIVE_CHUNK_RADIUS {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, "can only load nearby chunks", http.StatusForbidden)
			return
		}

		chunks := types.NewChunks(radius)

		for _, x := range makeRange(center.X-int64(radius), center.X+int64(radius)) {
			for _, y := range makeRange(center.Y-int64(radius), center.Y+int64(radius)) {
				for _, z := range makeRange(center.Z-int64(radius), center.Z+int64(radius)) {
					point := types.NewPoint(x, y, z)
					chunk, err := engine.World.Chunk(point)
					if err != nil {
						Log.Errorf("web:chunk: error=%v", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else if chunk == nil { // shouldn't be possible
						Log.Errorf(`web:chunk: error=No such Chunk "%s"`, point.ToString())
						http.Error(w, fmt.Sprintf(`No such Chunk "%s"`, point.ToString()), http.StatusNotFound)
						return
					}

					chunks.Chunks[point] = *chunk

					for i := 0; i < len(chunk.Entities); i++ {
						id := chunk.Entities[i]
						entity, err := engine.World.Entity(id)
						if err != nil {
							Log.Errorf("web:chunk: error=%v", err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						chunks.Entities[id] = *entity
					}

					for i := 0; i < len(chunk.Items); i++ {
						id := chunk.Items[i]
						item, err := engine.World.Item(id)
						if err != nil {
							Log.Errorf("web:chunk: error=%v", err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						chunks.Items[id] = *item
					}
				}
			}
		}

		if blob, err := chunks.Marshal(); err == nil {
			w.Header().Set("Content-Type", "application/protobuf")
			if _, err := w.Write(blob); err != nil {
				Log.Errorf("web:chunk: error=%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			Log.Errorf("web:chunk: error=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func worldAgeHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("web: worldAge")
		age, err := engine.World.Age()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if blob, err := age.Marshal(); err == nil {
			w.Header().Set("Content-Type", "application/protobuf")
			if _, err := w.Write(blob); err != nil {
				Log.Errorf("web:age: error=%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			Log.Errorf("web:age: error=%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func actHandler(engine *Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("web:act")
		if r.Body == nil {
			http.Error(w, "Must supply body", http.StatusBadRequest)
			return
		}

		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		action := types.Action{}
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

func makeRange(min, max int64) []int64 {
	a := make([]int64, max-min+1)
	for i := range a {
		a[i] = min + int64(i)
	}
	return a
}
