package engine

import (
	"bytes"
	"fmt"
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {
	g := Goblin(t)

	g.Describe("http api", func() {
		NAME := "arana"
		PASS := "murakami"

		auth, _ := (&types.Auth{Name: NAME, Password: []byte(PASS)}).Marshal()

		var engine *Engine
		var token string

		g.Before(func() {
			unique, err := uuid.NewV4()
			if err != nil {
				t.Fatal(err)
			}

			engine = &Engine{}
			if err := engine.Init("/tmp/savedir-huyilla-" + unique.String()); err != nil {
				t.Fatal(err)
			}
		})

		g.After(func() {
			if err := engine.World.WipeDatabase(); err != nil {
				t.Fatal(err)
			}
		})

		g.Describe("ping", func() {
			g.It("ping returns pong", func() {
				res := requesty("GET", "/ping", nil, engine, nil)

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("pong"))
			})
		})

		g.Describe("signup flow", func() {
			g.It("Signs up", func() {
				res := requesty("POST", "/auth/signup", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("Signup successful!"))
				g.Assert(res.Code).Equal(http.StatusOK)

				player, err := engine.World.Player(NAME)
				g.Assert(err).IsNil()
				g.Assert(player).IsNotNil()
				g.Assert(len(player.Token)).Equal(0)
				g.Assert(player.Name).Equal(NAME)
			})

			g.It("Logs in", func() {
				res := requesty("POST", "/auth/login", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(len(body) > 100).IsTrue(fmt.Sprintf(`Body was too short: "%s"`, body))
				g.Assert(body[0]).Equal(byte('e'))
				g.Assert(res.Code).Equal(http.StatusOK)

				token = string(body)

				player, err := engine.World.Player(NAME)
				g.Assert(err).IsNil()
				g.Assert(player).IsNotNil()
				g.Assert(len(player.Token) > 0).IsTrue("Player token is not set")
				g.Assert(player.Name).Equal(NAME)
			})

			g.It("Logs out", func() {
				res := requesty("POST", "/auth/logout", nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("Logout successful!"))

				player, err := engine.World.Player(NAME)
				g.Assert(err).IsNil()
				g.Assert(player).IsNotNil()
				g.Assert(len(player.Token)).Equal(0)
			})
		})

		g.Describe("get wold age", func() {
			g.It("gets world age", func() {
				res := requesty("GET", "/world/age", nil, engine, nil)

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("1"))
			})
		})

		g.Describe("get player", func() {
			g.It("Gets player info", func() {
				res := requesty("GET", "/world/player/"+NAME, nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				var entity types.Entity
				if err := entity.Unmarshal(body); err != nil {
					t.Fatal(err)
				}
				g.Assert(entity.PlayerName).Equal(NAME)
				g.Assert(entity.Id).NotEqual(0)
			})
		})

		g.Describe("get chunk", func() {
			g.It("in range", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				var chunk types.Chunk
				if err := chunk.Unmarshal(body); err != nil {
					t.Fatal(err)
				}
				g.Assert(len(chunk.Voxels)).Equal(constants.CHUNK_LENGTH)
			})

			g.It("out of range", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d", 0, constants.ACTIVE_CHUNK_RADIUS+1, 0), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusForbidden)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(string(body)).Equal("can only load nearby chunks\n")
			})
		})
	})
}

func requesty(method, url string, body io.Reader, engine *Engine, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res := httptest.NewRecorder()
	Router(engine).ServeHTTP(res, req)

	return res
}
