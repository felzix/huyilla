package engine

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_Chunk(t *testing.T) {
	h := &Engine{}
	h.Init("/tmp/huyilla")
	defer h.World.WipeDatabase()

	if _, err := h.World.Chunk(&types.Point{X: 0, Y: 0, Z: 0}); err != nil {
		t.Fatal(err)
	}

	chunk, err := h.World.Chunk(&types.Point{X: 0, Y: 0, Z: 0})
	if err != nil {
		t.Fatal(err)
	}

	expectedVoxelCount := C.CHUNK_SIZE * C.CHUNK_SIZE * C.CHUNK_SIZE
	if len(chunk.Voxels) != expectedVoxelCount {
		t.Errorf(`Was expected %d voxels but got %d`, expectedVoxelCount, len(chunk.Voxels))
	}
}
