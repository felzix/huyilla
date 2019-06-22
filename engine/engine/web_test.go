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
			g.Assert(err).IsNil()

			engine = &Engine{}
			err = engine.Init("/tmp/savedir-huyilla-" + unique.String())
			g.Assert(err).IsNil()
		})

		g.After(func() {
			err := engine.World.WipeDatabase()
			g.Assert(err).IsNil()
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
			g.It("Doesn't yet exist", func() {
				res := requesty("GET", "/auth/exists/"+NAME, nil, engine, nil)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("false"))
				g.Assert(res.Code).Equal(http.StatusOK)
			})

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

			g.It("Now exists", func() {
				res := requesty("GET", "/auth/exists/"+NAME, nil, engine, nil)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(body).Equal([]byte("true"))
				g.Assert(res.Code).Equal(http.StatusOK)
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

		g.Describe("get world age", func() {
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

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				g.Assert(err).IsNil()
				g.Assert(len(chunks.Chunks)).Equal(1)
				g.Assert(chunks.Chunks[0]).IsNotNil()
				g.Assert(len(chunks.Chunks[0].Voxels)).Equal(constants.CHUNK_LENGTH)
			})

			g.It("in range, radius=1", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=1", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				g.Assert(err).IsNil()
				g.Assert(len(chunks.Chunks)).Equal(27)
				g.Assert(chunks.Chunks[0]).IsNotNil()
				g.Assert(chunks.Chunks[26]).IsNotNil()
				g.Assert(len(chunks.Chunks[0].Voxels)).Equal(constants.CHUNK_LENGTH)
			})

			g.It("barely in range, radius is max", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 0, 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				g.Assert(err).IsNil()
				g.Assert(len(chunks.Chunks)).Equal(7 * 7 * 7)
				g.Assert(chunks.Chunks[0]).IsNotNil()
				g.Assert(chunks.Chunks[7*7*7-1]).IsNotNil()
				g.Assert(len(chunks.Chunks[0].Voxels)).Equal(constants.CHUNK_LENGTH)
			})

			// The idea is that the initial readtime is so large that two of them don't fit into the per-test duration set on the commandline (15s)
			g.It("database caching works", func() {
				// fill cache
				requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=3", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				// read from cache
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=3", 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()

				var chunks types.Chunks
				err = chunks.Unmarshal(body)
				g.Assert(err).IsNil()
				g.Assert(len(chunks.Chunks)).Equal(7 * 7 * 7)
				g.Assert(chunks.Chunks[0]).IsNotNil()
				g.Assert(chunks.Chunks[7*7*7-1]).IsNotNil()
				g.Assert(len(chunks.Chunks[0].Voxels)).Equal(constants.CHUNK_LENGTH)
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

			g.It("out of range: too big", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 0, 0, 0, constants.ACTIVE_CHUNK_RADIUS+1), nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token":       token,
				})

				g.Assert(res.Code).Equal(http.StatusForbidden)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).IsNil()
				g.Assert(string(body)).Equal("can only load nearby chunks\n")
			})

			g.It("out of range: max size and shifted", func() {
				res := requesty("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", 1, 0, 0, constants.ACTIVE_CHUNK_RADIUS), nil, engine, map[string]string{
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
