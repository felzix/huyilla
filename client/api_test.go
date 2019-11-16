package client

import (
	"context"
	"fmt"
	"github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/engine/engine"
	"github.com/felzix/huyilla/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = Describe("HTTP API", func() {
	const PORT = 8085

	api := NewAPI(fmt.Sprintf("http://localhost:%d", PORT), "arana", "murakami")

	var huyilla *engine.Engine
	var webServerError chan error
	var server *http.Server

	BeforeSuite(func() {
		engine, err := engine.NewEngine(constants.SEED, engine.NewLakeWorldGenerator(3), engine.NewMemoryDatabase())
		huyilla = engine
		Expect(err).To(BeNil())
		server = huyilla.Serve(PORT, webServerError)
	})

	AfterSuite(func() {
		err := server.Shutdown(context.TODO())
		Expect(err).To(BeNil())
	})

	Describe("ping", func() {
		It("ping returns pong", func() {
			response, err := api.Ping()
			Expect(err).To(BeNil())
			Expect(response).To(Equal("pong"))
		})
	})

	Describe("signup flow", func() {
		It("checks if user already exists", func() {
			exists, err := api.UserExists()
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())
		})

		It("signs up", func() {
			err := api.Signup()
			Expect(err).To(BeNil())

			var player *types.Player
			Eventually(func() *types.Player {
				player, _ = huyilla.World.Player(api.Username)
				return player
			}).ShouldNot(BeNil())

			Expect(len(player.Token)).To(Equal(0))
			Expect(player.Name).To(Equal(api.Username))

			entity, err := huyilla.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())

			Eventually(func() *types.Player {
				player, _ = huyilla.World.Player(api.Username)
				return player
			})
			Expect(len(player.Token)).To(Equal(0))
			Expect(player.Name).To(Equal(api.Username))

			entity, err = huyilla.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())
		})

		It("checks if user now exists", func() {
			exists, err := api.UserExists()
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})

		It("logs in", func() {
			err := api.Login()
			Expect(err).To(BeNil())

			var player *types.Player
			Eventually(func() *types.Player {
				player, _ = huyilla.World.Player(api.Username)
				return player
			}).ShouldNot(BeNil())

			Expect(len(player.Token) > 0).To(BeTrue(), "Player token is not set")
			Expect(player.Name).To(Equal(api.Username))
		})

		It("logs out", func() {
			err := api.Logout()
			Expect(err).To(BeNil())

			var player *types.Player
			Eventually(func() *types.Player {
				player, _ = huyilla.World.Player(api.Username)
				if len(player.Token) != 0 {
					return nil
				}
				return player
			}).ShouldNot(BeNil())

			// Can log back in
			err = api.Login()
			Expect(err).To(BeNil())

			Eventually(func() *types.Player {
				player, _ = huyilla.World.Player(api.Username)
				return player
			}).ShouldNot(BeNil())

			Expect(len(player.Token) > 0).To(BeTrue(),"Player token is not set")
			Expect(player.Name).To(Equal(api.Username))
		})

	})

	Describe("world getting", func() {
		It("gets world age", func() {
			age, err := api.GetWorldAge()
			Expect(err).To(BeNil())
			Expect(age).To(Equal(types.Age(1)))
		})

		It("gets player", func() {
			player, err := api.GetPlayer(api.Username)
			Expect(err).To(BeNil())
			Expect(player).ToNot(BeNil())
			Expect(player.PlayerName).To(Equal(api.Username))
			Expect(player.Id).ToNot(Equal(0))
		})

		It("gets one chunk in range", func() {
			point := types.NewPoint(0, 0, constants.ACTIVE_CHUNK_RADIUS)
			chunks, err := api.GetChunks(types.NewPoint(point.X, point.Y, point.Z), 0)
			Expect(err).To(BeNil())
			Expect(len(chunks.Chunks)).To(Equal(1))
			Expect(len(chunks.Chunks[point].Voxels)).To(Equal(constants.CHUNK_LENGTH))
		})

		It("cannot get a chunk out of range", func() {
			_, err := api.GetChunks(types.NewPoint(0, 0, constants.ACTIVE_CHUNK_RADIUS+1), 0)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("GetChunk failure: Expected status 200 but got 403. can only load nearby chunks\n"))
		})
	})

	Describe("action issuing", func() {
		It("issues an action", func() {
			player, err := api.GetPlayer(api.Username)
			Expect(err).To(BeNil())

			originalY := player.Location.Voxel.Y
			player.Location.Voxel.Y++
			err = api.IssueMoveAction(&types.Player{EntityId: player.Id}, player.Location)
			Expect(err).To(BeNil())

			err = huyilla.Tick()
			Expect(err).To(BeNil())

			player, err = api.GetPlayer(api.Username)
			Expect(err).To(BeNil())
			Expect(player.Location.Voxel.Y).To(Equal(originalY + 1))
		})
	})
})

