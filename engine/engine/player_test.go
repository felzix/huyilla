package engine

import (
	uuid "github.com/satori/go.uuid"
	"testing"
	. "github.com/felzix/goblin"
)

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("Content Test", func() {
		NAME := "felzix"
		PASS := "murakami"
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
			err := h.SignUp(NAME, PASS)
			g.Assert(err).IsNil()
			_, err = h.LogIn(NAME, PASS)
			g.Assert(err).IsNil()

			player, err := h.World.Player(NAME)
			g.Assert(err).IsNil()
			g.Assert(player).IsNotNil()
			g.Assert(player.Name).Equal(NAME)
		})
	})
}
