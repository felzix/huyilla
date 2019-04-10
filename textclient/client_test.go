package main

import (
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/content"
	"testing"
)

func TestClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("Client Test", func() {
		g.Before(func() {
			content.PopulateContentNameMaps()
		})

		g.It("air voxel", func() {
			g.Assert(voxelToRune(uint64(0))).Equal(' ')
		})

		g.It("barren earth voxel", func() {
			g.Assert(voxelToRune(uint64(1))).Equal('.')
		})
	})
}
