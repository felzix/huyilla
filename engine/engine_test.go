package main

import (
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_ActiveChunkRadius(t *testing.T) {
	h := &Engine{}
	h.Init(&types.Config{})

	if err := h.SignUp("felzix", "PASS"); err != nil {
		t.Fatal(err)
	}
	player, err := h.LogIn("felzix", "PASS")
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Tick(); err != nil {
		t.Fatal(err)
	}

	chunk, err := h.GetChunk(player.Entity.Location.Chunk)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunk.Entities) != 1 {
		t.Errorf("Expected 1 entity in chunk but there were %d", len(chunk.Entities))
	}

	// active range in the positive direction
	edge := clonePoint(player.Entity.Location.Chunk)
	edge.X += 3
	chunk, err = h.GetChunk(edge)
	if chunk == nil {
		t.Error("Chunk within player's range should exist but it does not.")
	}
	if len(chunk.Entities) != 0 {
		t.Errorf("Expected 0 entities in chunk but there were %d", len(chunk.Entities))
	}

	beyond := clonePoint(player.Entity.Location.Chunk)
	beyond.X += 4
	chunk, err = h.GetChunk(beyond)
	if err == nil { // note that this is "==" not "!="
		t.Error("Chunk beyond player's range exists when it should not.")
	}

	// active range in the negative direction
	edge = clonePoint(player.Entity.Location.Chunk)
	edge.X -= 3
	chunk, err = h.GetChunk(edge)
	if err != nil {
		t.Fatal(err)
	}
	if chunk == nil {
		t.Error("Chunk within player's range should exist but it does not.")
	}

	beyond = clonePoint(player.Entity.Location.Chunk)
	beyond.X -= 4
	chunk, err = h.GetChunk(beyond)
	if err == nil { // note that this is "==" not "!="
		t.Error("Chunk beyond player's range exists when it should not.")
	}
}
