package engine

import (
	"testing"
)

func TestHuyilla_ActiveChunkRadius(t *testing.T) {
	h := &Engine{}
	h.Init("/tmp/huyilla")
	defer h.World.WipeDatabase()

	if err := h.SignUp("felzix", "PASS"); err != nil {
		t.Fatal(err)
	}
	_, err := h.LogIn("felzix", "PASS")
	if err != nil {
		t.Fatal(err)
	}

	if err := h.Tick(); err != nil {
		t.Fatal(err)
	}

	player, err := h.World.Player("felzix")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	entity, err := h.World.Entity(player.EntityId)
	if entity == nil {
		t.Fatalf("Entity %d should exist but doesn't", player.EntityId)
	} else if err != nil {
		t.Fatal(err)
	}

	chunk, err := h.World.Chunk(entity.Location.Chunk)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunk.Entities) != 1 {
		t.Errorf("Expected 1 entity in chunk but there were %d", len(chunk.Entities))
	}

	// active range in the positive direction
	edge := clonePoint(entity.Location.Chunk)
	edge.X += 3
	chunk, err = h.World.Chunk(edge)
	if chunk == nil {
		t.Error("Chunk within player's range should exist but it does not.")
	}
	if len(chunk.Entities) != 0 {
		t.Errorf("Expected 0 entities in chunk but there were %d", len(chunk.Entities))
	}

	beyond := clonePoint(entity.Location.Chunk)
	beyond.X += 4
	chunk, err = h.World.OnlyGetChunk(beyond)
	if chunk != nil {
		t.Error("Chunk beyond player's range exists when it should not.")
	}

	// active range in the negative direction
	edge = clonePoint(entity.Location.Chunk)
	edge.X -= 3
	chunk, err = h.World.Chunk(edge)
	if err != nil {
		t.Fatal(err)
	}
	if chunk == nil {
		t.Error("Chunk within player's range should exist but it does not.")
	}

	beyond = clonePoint(entity.Location.Chunk)
	beyond.X -= 4
	chunk, err = h.World.OnlyGetChunk(beyond)
	if chunk != nil {
		t.Error("Chunk beyond player's range exists when it should not.")
	}
}
