package main

import (
	"testing"
)

func TestHuyilla_Entity(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	if err := h.SignUp("felzix", "PASS"); err != nil {
		t.Fatal(err)
	}

	player, err := h.World.Player("felzix")
	if err != nil {
		t.Fatal(err)
	}

	entity, err := h.World.Entity(player.EntityId)
	if entity == nil {
		t.Fatalf("Entity %d should exist but doesn't", player.EntityId)
	} else if err != nil {
		t.Fatal(err)
	}

	if entity.Type != entity.Type {
		t.Errorf(`GetPlayer and GetEntity returned different entities: "%v" != "%v"`, entity, entity)
	}
}
