package main

import (
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("Client Test", func() {
		g.It("air voxel", func() {
			air := types.ExpandedVoxel{
				Form:     0, // cube
				Material: 0, // air
			}.Compress()

			g.Assert(voxelToRune(air)).Equal(' ')
		})

		g.It("barren earth voxel", func() {
			dirt := types.ExpandedVoxel{
				Form:     0,   // cube
				Material: 100, // dirt
			}.Compress()
			g.Assert(voxelToRune(dirt)).Equal('.')
		})
	})
}
