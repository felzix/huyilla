package main

import (
	"context"
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/constants"
	engine2 "github.com/felzix/huyilla/engine/engine"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestAPI(t *testing.T) {
	g := Goblin(t)

	g.Describe("http api", func() {
		api := NewAPI("http://localhost:8080", "arana", "murakami")

		var engine *engine2.Engine
		var webServerError chan error
		var server *http.Server

		g.Before(func() {
			unique, err := uuid.NewV4()
			if err != nil {
				t.Fatal(err)
			}

			engine = &engine2.Engine{}
			if err := engine.Init("/tmp/savedir-huyilla-" + unique.String()); err != nil {
				t.Fatal(err)
			}
			server = engine.Serve(webServerError)
		})

		g.After(func() {
			if err := engine.World.WipeDatabase(); err != nil {
				t.Fatal(err)
			}
			if err := server.Shutdown(context.TODO()); err != nil {
				t.Fatal(err)
			}
		})

		g.Describe("ping", func() {
			g.It("ping returns pong", func() {
				response, err := api.Ping()
				g.Assert(err).Equal(nil)
				g.Assert(response).Equal("pong")
			})
		})

		g.Describe("signup flow", func() {
			g.It("checks if user already exists", func() {
				exists, err := api.UserExists()
				g.Assert(err).Equal(nil)
				g.Assert(exists).Equal(false)
			})

			g.It("signs up", func() {
				err := api.Signup()
				g.Assert(err).Equal(nil)

				g.Poll(5, 200, func() bool {
					player, err := engine.World.Player(api.Username)
					g.Assert(err).Equal(nil)

					if player == nil {
						return false
					}

					g.Assert(len(player.Token)).Equal(0)
					g.Assert(player.Name).Equal(api.Username)

					entity, err := engine.World.Entity(player.EntityId)
					g.Assert(err).Equal(nil)
					g.Assert(entity).NotEqual(nil)

					return true
				})
			})

			g.It("checks if user now exists", func() {
				exists, err := api.UserExists()
				g.Assert(err).Equal(nil)
				g.Assert(exists).Equal(true)
			})


			g.It("logs in", func() {
				err := api.Login()
				g.Assert(err).Equal(nil)

				g.Poll(5, 200, func() bool {
					player, err := engine.World.Player(api.Username)
					g.Assert(err).Equal(nil)

					if player == nil {
						return false
					}

					g.Assert(len(player.Token) > 0).IsTrue("Player token is not set")
					g.Assert(player.Name).Equal(api.Username)

					return true
				})
			})

			g.It("logs out", func() {
				err := api.Logout()
				g.Assert(err).Equal(nil)

				var player *types.Player
				g.Poll(5, 200, func() bool {
					player, err = engine.World.Player(api.Username)

					if player == nil || len(player.Token) != 0 {
						return false
					}

					return true
				})

				// Can log back in
				err = api.Login()
				g.Assert(err).Equal(nil)

				g.Poll(5, 200, func() bool {
					player, err = engine.World.Player(api.Username)
					g.Assert(err).Equal(nil)

					if player == nil {
						return false
					}
					g.Assert(len(player.Token) > 0).IsTrue("Player token is not set")
					g.Assert(player.Name).Equal(api.Username)

					return true
				})
			})

		})

		g.Describe("world getting", func() {
			g.It("gets player", func() {
				player, err := api.GetPlayer(api.Username)
				g.Assert(err).Equal(nil)
				g.Assert(player).NotEqual(nil)
				g.Assert(player.Name).Equal(api.Username)
				g.Assert(player.EntityId).NotEqual(0)
				g.Assert(player.Password).Equal([]byte(nil)) // don't transmit password
				g.Assert(player.Spawn).Equal((*types.AbsolutePoint)(nil)) // don't transmit spawn point
			})

			g.It("gets chunk", func() {
				chunk, err := api.GetChunk(engine2.NewPoint(0, 0, constants.ACTIVE_CHUNK_RADIUS))
				g.Assert(err).Equal(nil)
				g.Assert(len(chunk.Voxels)).Equal(constants.CHUNK_LENGTH)
			})

			g.It("cannot get chunk out of range", func() {
				_, err := api.GetChunk(engine2.NewPoint(0, 0, constants.ACTIVE_CHUNK_RADIUS + 1))
				g.Assert(err).NotEqual(nil)
				g.Assert(err.Error()).Equal("GetChunk failure: Expected status 200 but got 403. can only load nearby chunks\n")
			})
		})
	})
}
