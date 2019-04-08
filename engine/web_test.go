package main

import (
	"bytes"
	"fmt"
	"github.com/felzix/huyilla/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/felzix/goblin"
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
			engine = &Engine{}
			if err := engine.Init(); err != nil {
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
				g.Assert(err).Equal(nil)
				g.Assert(body).Equal([]byte("pong"))
			})
		})

		g.Describe("signup flow", func() {
			g.It("Signs up", func() {
				res := requesty("POST", "/auth/signup", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).Equal(nil)
				g.Assert(body).Equal([]byte("Signup successful!"))
				g.Assert(res.Code).Equal(http.StatusOK)

				player, err := engine.World.Player(NAME)
				g.Assert(err).Equal(nil)
				g.Assert(player).NotEqual(nil)
				g.Assert(len(player.Token)).Equal(0)
				g.Assert(player.Name).Equal(NAME)
			})

			g.It("Logs in", func() {
				res := requesty("POST", "/auth/login", bytes.NewReader(auth), engine, map[string]string{
					"contentType": "application/protobuf",
				})

				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).Equal(nil)
				g.Assert(len(body) > 100).IsTrue(fmt.Sprintf(`Body was too short: "%s"`, body))
				g.Assert(body[0]).Equal(byte('e'))
				g.Assert(res.Code).Equal(http.StatusOK)

				token = string(body)

				player, err := engine.World.Player(NAME)
				g.Assert(err).Equal(nil)
				g.Assert(player).NotEqual(nil)
				g.Assert(len(player.Token) > 0).IsTrue("Player token is not set")
				g.Assert(player.Name).Equal(NAME)
			})

			g.It("Logs out", func() {
				res := requesty("POST", "/auth/logout", nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token": token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).Equal(nil)
				g.Assert(body).Equal([]byte("Logout successful!"))

				player, err := engine.World.Player(NAME)
				g.Assert(err).Equal(nil)
				g.Assert(player).NotEqual(nil)
				g.Assert(len(player.Token)).Equal(0)
			})
		})

		g.Describe("get player", func() {
			g.It("Gets player info", func() {
				res := requesty("GET", "/world/player/"+NAME, nil, engine, map[string]string{
					"contentType": "application/protobuf",
					"token": token,
				})

				g.Assert(res.Code).Equal(http.StatusOK)
				body, err := ioutil.ReadAll(res.Body)
				g.Assert(err).Equal(nil)
				var player types.Player
				if err := player.Unmarshal(body); err != nil {
					t.Fatal(err)
				}
				g.Assert(player.Name).Equal(NAME)
				g.Assert(player.EntityId).NotEqual(0)
				g.Assert(player.Password).Equal([]byte(nil))
				g.Assert(player.Token).Equal("")
				g.Assert(player.Spawn).Equal((*types.AbsolutePoint)(nil))
			})

		})
	})
}

func requesty (method, url string, body io.Reader, engine *Engine, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res := httptest.NewRecorder()
	Router(engine).ServeHTTP(res, req)

	return res
}
