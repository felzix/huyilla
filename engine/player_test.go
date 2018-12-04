package main

import (
	"testing"
)

func TestHuyilla_Players(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	NAME := "felzix"
	PASS := "murakami"

	h.SignUp(NAME, PASS)
	h.LogIn(NAME, PASS)

	player, err := h.World.Player(NAME)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.Name != "felzix" {
		t.Errorf(`Player name was "%v" instead of "felzix"`, player.Name)
	}
}
