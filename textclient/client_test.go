package main

import (
	"github.com/felzix/huyilla/content"
	"testing"
)

func TestClient_voxelToRune(t *testing.T) {
	content.PopulateContentNameMaps()

	if runic := voxelToRune(uint64(0)); runic != ' ' {
		t.Errorf(`Air voxel should have rune ' ' but has '%v'`, runic)
	}

	if runic := voxelToRune(uint64(1)); runic != '.' {
		t.Errorf(`Barren earth voxel should have rune ' ' but has '%v'`, runic)
	}
}
