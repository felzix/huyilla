package engine

import (
	. "github.com/felzix/goblin"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestContent(t *testing.T) {
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
			content := h.GetContent()
			g.Assert(content.E[0].Name).Equal("human")
		})
	})
}
