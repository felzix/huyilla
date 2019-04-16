package engine

import (
	"fmt"
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestEngine(t *testing.T) {
	g := Goblin(t)
	g.Describe("Engine Test", func() {
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

		g.It("ticks", func() {
			err := h.SignUp("felzix", "PASS")
			g.Assert(err).IsNil()

			_, err = h.LogIn("felzix", "PASS")
			g.Assert(err).IsNil()

			err = h.Tick()
			g.Assert(err).IsNil()

			content := h.GetContent()
			g.Assert(content.E[0].Name).Equal("human")

			player, err := h.World.Player("felzix")
			g.Assert(err).IsNil()

			entity, err := h.World.Entity(player.EntityId)
			g.Assert(err).IsNil()
			g.Assert(entity).IsNotNil()

			chunk, err := h.World.Chunk(entity.Location.Chunk)
			g.Assert(err).IsNil()
			g.Assert(len(chunk.Entities)).Equal(1)

			// active range in the positive direction
			edge := types.ClonePoint(entity.Location.Chunk)
			edge.X += 3
			chunk, err = h.World.Chunk(edge)
			g.Assert(chunk).IsNotNil() // chunk within player's range should exist
			g.Assert(len(chunk.Entities)).Equal(0)

			beyond := types.ClonePoint(entity.Location.Chunk)
			beyond.X += 4
			chunk, err = h.World.OnlyGetChunk(beyond)

			fmt.Printf("%T , %v", chunk, chunk)

			g.Assert(chunk).IsNil() // Chunk beyond player's range exists

			// active range in the negative direction
			edge = types.ClonePoint(entity.Location.Chunk)
			edge.X -= 3
			chunk, err = h.World.Chunk(edge)
			g.Assert(err).IsNil()
			g.Assert(chunk).IsNotNil() // Chunk within player's range should exist

			beyond = types.ClonePoint(entity.Location.Chunk)
			beyond.X -= 4
			chunk, err = h.World.OnlyGetChunk(beyond)
			g.Assert(err).IsNil()
			g.Assert(chunk).IsNil() // Chunk beyond player's range exists
		})
	})
}