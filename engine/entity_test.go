package main

import (
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_Entity(t *testing.T) {
	h := &Engine{}
	h.Init(&types.Config{})

	if err := h.SignUp("felzix", "PASS"); err != nil {
		t.Fatal(err)
	}

	player, err := h.GetPlayer("felzix")
	if err != nil {
		t.Fatal(err)
	}

	entity := h.Entities[player.Player.EntityId]

	if entity.Type != player.Entity.Type {
		t.Errorf(`GetPlayer and GetEntity returned different entities: "%v" != "%v"`, entity, player.Entity)
	}
}
