package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"testing"
	. "github.com/felzix/goblin"
)

func TestChunk(t *testing.T) {
	g := Goblin(t)
	g.Describe("Chunk test", func() {
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

		g.It("generates a chunk", func() {
			chunk, err := h.World.Chunk(&types.Point{X: 0, Y: 0, Z: 0})
			g.Assert(err).IsNil()

			expectedVoxelCount := C.CHUNK_SIZE * C.CHUNK_SIZE * C.CHUNK_SIZE
			g.Assert(len(chunk.Voxels)).Equal(expectedVoxelCount)
		})
	})
}
