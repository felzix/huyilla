package main

import (
	"testing"
)

func TestHuyilla_SignUp(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	if err := h.SignUp("felzix", "PASS"); err != nil {
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

	if player.Name != entity.PlayerName {
		t.Errorf(
			`Player's entity has wrong name: player="%v", entity="%v"`,
			player.Name, entity.PlayerName)
	}
}

func TestHuyilla_Login(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	if err := h.SignUp("felzix", "PASS"); err != nil {
		t.Fatal(err)
	}

	details, err := h.LogIn("felzix", "PASS")
	if err != nil {
		t.Fatal(err)
	}

	if details.Player.Spawn == nil {
		t.Error("Player should have gotten a default spawn but did not")
	}

	if details.Entity.Location == nil {
		t.Fatal("Player entity should have been created but it was not (no location)")
	}

	chunk, err := h.World.Chunk(details.Entity.Location.Chunk)
	if err != nil {
		t.Fatal(err)
	}

	entityIsPresent := false
	for i := 0; i < len(chunk.Entities); i++ {
		entity := chunk.Entities[i]
		if entity == details.Entity.Id {
			entityIsPresent = true
		}
	}
	if !entityIsPresent {
		t.Errorf(`Expected entity at chunk (%d,%d,%d) but it was not there`,
			details.Entity.Location.Chunk.X,
			details.Entity.Location.Chunk.Y,
			details.Entity.Location.Chunk.Z)
	}
}

func TestHuyilla_LoginNegative(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	_, err := h.LogIn("felzix", "PASS")
	if err == nil {
		t.Fatal("Logging in before signup should throw an error but didn't")
	} else if err.Error() != `No such player "felzix"` {
		t.Errorf(`Wrong error. Got "%v"`, err)
	}
}
