package engine

import (
	"fmt"
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("Auth", func() {
		var h *Engine

		g.BeforeEach(func() {
			unique, err := uuid.NewV4()
			if err != nil {
				t.Fatal(err)
			}

			h = &Engine{}
			if err := h.Init("/tmp/savedir-huyilla-" + unique.String()); err != nil {
				t.Fatal(err)
			}
		})

		g.AfterEach(func() {
			if h == nil || h.World == nil {
				return
			}
			if err := h.World.WipeDatabase(); err != nil {
				t.Fatal(err)
			}
		})

		g.It("signs up", func() {
			err := h.SignUp("felzix", "PASS")
			g.Assert(err).Equal(nil)

			player, err := h.World.Player("felzix")
			g.Assert(err).Equal(nil)
			g.Assert(player).NotEqual(nil)

			entity, err := h.World.Entity(player.EntityId)
			g.Assert(err).Equal(nil)
			g.Assert(entity).NotEqual(nil)

			g.Assert(player.Name).Equal(entity.PlayerName)
		})

		g.It("logs in", func() {
			err := h.SignUp("felzix", "PASS")
			g.Assert(err).Equal(nil)

			token, err := h.LogIn("felzix", "PASS")
			g.Assert(err).Equal(nil)
			g.Assert(len(token) > 100).IsTrue("token is not set or is set incorrectly")
			g.Assert(token[0]).Equal(byte('e'))

			player, err := h.World.Player("felzix")
			g.Assert(err).Equal(nil)
			g.Assert(player).NotEqual(nil)
			g.Assert(player.Token).Equal(token)

			entity, err := h.World.Entity(player.EntityId)
			g.Assert(err).Equal(nil)
			g.Assert(entity).NotEqual(nil)

			chunk, err := h.World.Chunk(entity.Location.Chunk)
			g.Assert(err).Equal(nil)

			entityIsPresent := false
			for i := 0; i < len(chunk.Entities); i++ {
				entity := chunk.Entities[i]
				if entity == player.EntityId {
					entityIsPresent = true
				}
			}
			g.Assert(entityIsPresent).IsTrue(fmt.Sprintf(
				`Expected entity at chunk (%d,%d,%d) but it was not there`,
				entity.Location.Chunk.X,
				entity.Location.Chunk.Y,
				entity.Location.Chunk.Z))
		})

		g.It("logs out", func() {
			err := h.SignUp("felzix", "PASS")
			g.Assert(err).Equal(nil)

			_, err = h.LogIn("felzix", "PASS")
			g.Assert(err).Equal(nil)

			err = h.LogOut("felzix")
			g.Assert(err).Equal(nil)

			player, err := h.World.Player("felzix")
			g.Assert(err).Equal(nil)
			g.Assert(player).NotEqual(nil)
			g.Assert(player.Token).Equal("")

			entity, err := h.World.Entity(player.EntityId)
			g.Assert(err).Equal(nil)
			g.Assert(entity).NotEqual(nil)

			chunk, err := h.World.Chunk(entity.Location.Chunk)
			g.Assert(err).Equal(nil)

			entityIsPresent := false
			for i := 0; i < len(chunk.Entities); i++ {
				entity := chunk.Entities[i]
				if entity == player.EntityId {
					entityIsPresent = true
				}
			}
			g.Assert(entityIsPresent).IsFalse(fmt.Sprintf(
				`Expected entity to NOT be at chunk (%d,%d,%d) but it was there.`,
				entity.Location.Chunk.X,
				entity.Location.Chunk.Y,
				entity.Location.Chunk.Z))
		})

		g.It("fails to log in", func() {
			token, err := h.LogIn("felzix", "PASS")
			g.Assert(token).Equal("")
			g.Assert(err).NotEqual(nil)
			g.Assert(err.Error()).Equal(`No such player "felzix"`)

			player, err := h.World.Player("felzix")
			g.Assert(err).Equal(nil)
			g.Assert(player).Equal((*types.Player)(nil))
		})

		g.It("token test", func() {
			SECRET := []byte("secret")
			NAME := "camian"
			EXPIRY := time.Now().Add(time.Hour * 24).Unix()

			token, err := makeToken(SECRET, NAME, EXPIRY)
			g.Assert(err).Equal(nil)

			name, tokenId, expiry, err := readToken(SECRET, token)
			g.Assert(err).Equal(nil)

			g.Assert(name).Equal(NAME)
			g.Assert(expiry).Equal(EXPIRY)
			g.Assert(len(tokenId) > 0).IsTrue("Expected a token id but it's empty")
		})
	})
}
