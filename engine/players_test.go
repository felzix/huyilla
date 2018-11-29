package main

import (
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_Players(t *testing.T) {
	h := &Engine{}
	h.Init(&types.Config{})

	NAME := "felzix"
	PASS := "murakami"

	h.SignUp(NAME, PASS)
	h.LogIn(NAME, PASS)

	players := h.GetPlayerList()

	if len(players) != 1 {
		t.Errorf(`Error: Should be one player but there isn't: "%v"`, players)
	}

	player, err := h.GetPlayer(NAME)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.Player.Name != "felzix" {
		t.Errorf(`Player name was "%v" instead of "felzix"`, player.Player.Name)
	}
}
