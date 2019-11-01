package engine

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = Describe("Auth", func() {
		var h *Engine

		BeforeEach(func() {
			engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
			h = engine
			Expect(err).To(BeNil())
		})

		It("signs up", func() {
			err := h.SignUp("felzix", "PASS")
			Expect(err).To(BeNil())

			player, err := h.World.Player("felzix")
			Expect(err).To(BeNil())
			Expect(player).ToNot(BeNil())

			entity, err := h.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())

			Expect(player.Name).To(Equal(entity.PlayerName))
		})

		It("logs in", func() {
			err := h.SignUp("felzix", "PASS")
			Expect(err).To(BeNil())

			token, err := h.LogIn("felzix", "PASS")
			Expect(err).To(BeNil())
			Expect(len(token) > 100).To(BeTrue(), "token is not set or is set incorrectly")
			Expect(token[0]).To(Equal(byte('e')))

			player, err := h.World.Player("felzix")
			Expect(err).To(BeNil())
			Expect(player).ToNot(BeNil())
			Expect(player.Token).To(Equal(token))

			entity, err := h.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())

			chunk, err := h.World.Chunk(entity.Location.Chunk)
			Expect(err).To(BeNil())

			entityIsPresent := false
			for i := 0; i < len(chunk.Entities); i++ {
				entity := chunk.Entities[i]
				if entity == player.EntityId {
					entityIsPresent = true
				}
			}
			Expect(entityIsPresent).To(BeTrue(), fmt.Sprintf(
				`Expected entity at chunk (%d,%d,%d) but it was not there`,
				entity.Location.Chunk.X,
				entity.Location.Chunk.Y,
				entity.Location.Chunk.Z))
		})

		It("logs out", func() {
			err := h.SignUp("felzix", "PASS")
			Expect(err).To(BeNil())

			_, err = h.LogIn("felzix", "PASS")
			Expect(err).To(BeNil())

			err = h.LogOut("felzix")
			Expect(err).To(BeNil())

			player, err := h.World.Player("felzix")
			Expect(err).To(BeNil())
			Expect(player).ToNot(BeNil())
			Expect(player.Token).To(Equal(""))

			entity, err := h.World.Entity(player.EntityId)
			Expect(err).To(BeNil())
			Expect(entity).ToNot(BeNil())

			chunk, err := h.World.Chunk(entity.Location.Chunk)
			Expect(err).To(BeNil())

			entityIsPresent := false
			for i := 0; i < len(chunk.Entities); i++ {
				entity := chunk.Entities[i]
				if entity == player.EntityId {
					entityIsPresent = true
				}
			}
			Expect(entityIsPresent).To(BeFalse(), fmt.Sprintf(
				`Expected entity to NOT be at chunk (%d,%d,%d) but it was there.`,
				entity.Location.Chunk.X,
				entity.Location.Chunk.Y,
				entity.Location.Chunk.Z))
		})

		It("fails to log in", func() {
			token, err := h.LogIn("felzix", "PASS")
			Expect(token).To(Equal(""))
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(`No such player "felzix"`))

			player, err := h.World.Player("felzix")
			Expect(err).To(BeNil())
			Expect(player).To(Equal((*types.Player)(nil)))
		})

		It("token test", func() {
			SECRET := []byte("secret")
			NAME := "camian"
			EXPIRY := time.Now().Add(time.Hour * 24).Unix()

			token, err := makeToken(SECRET, NAME, EXPIRY)
			Expect(err).To(BeNil())

			name, tokenId, expiry, err := readToken(SECRET, token)
			Expect(err).To(BeNil())

			Expect(name).To(Equal(NAME))
			Expect(expiry).To(Equal(EXPIRY))
			Expect(len(tokenId) > 0).To(BeTrue(), "Expected a token id but it's empty")
		})
	})
