package engine

import (
	. "github.com/felzix/goblin"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestEntity(t *testing.T) {
	g := Goblin(t)
	g.Describe("Content Test", func() {
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

		g.It("loads human type", func() {
			err := h.SignUp("felzix", "PASS")
			g.Assert(err).IsNil()
			_, err = h.LogIn("felzix", "PASS")
			g.Assert(err).IsNil()

			player, err := h.World.Player("felzix")
			g.Assert(err).IsNil()
			g.Assert(player).IsNotNil()

			entity, err := h.World.Entity(player.EntityId)
			g.Assert(err).IsNil()
			g.Assert(entity).IsNotNil()
		})
	})
}
