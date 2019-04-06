package main

import (
	"testing"
	"time"
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

	if err := h.SignUp("allomance", "PASS"); err != nil {
		t.Fatal(err)
	}

	token, err := h.LogIn("allomance", "PASS")
	if err != nil {
		t.Fatal(err)
	}

	if len(token) < 100 || token[0] != 'e' {
		t.Errorf(`Bad token. Token="%s"`, token)
	}

	player, err := h.World.Player("allomance")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(player.Token) == 0 {
		t.Error(`Player should have been given a token`)
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

	entityIsPresent := false
	for i := 0; i < len(chunk.Entities); i++ {
		entity := chunk.Entities[i]
		if entity == player.EntityId {
			entityIsPresent = true
		}
	}
	if !entityIsPresent {
		t.Errorf(`Expected entity at chunk (%d,%d,%d) but it was not there`,
			entity.Location.Chunk.X,
			entity.Location.Chunk.Y,
			entity.Location.Chunk.Z)
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


func TestHuyilla_Logout(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	if err := h.SignUp("allomance", "PASS"); err != nil {
		t.Fatal(err)
	}

	_, err := h.LogIn("allomance", "PASS")
	if err != nil {
		t.Fatal(err)
	}

	if err := h.LogOut("allomance"); err != nil {
		t.Fatal(err)
	}

	player, err := h.World.Player("allomance")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(player.Token) > 0 {
		t.Error(`Player should not have a token`)
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

	entityIsPresent := false
	for i := 0; i < len(chunk.Entities); i++ {
		entity := chunk.Entities[i]
		if entity == player.EntityId {
			entityIsPresent = true
		}
	}
	if entityIsPresent {
		t.Errorf(`Unexpected entity at chunk (%d,%d,%d). It should have been removed on logout`,
			entity.Location.Chunk.X,
			entity.Location.Chunk.Y,
			entity.Location.Chunk.Z)
	}
}

func TestHuyilla_token(t *testing.T) {
	SECRET := []byte("secret")
	NAME := "camian"
	EXPIRY := time.Now().Add(time.Hour * 24).Unix()

	token, err := makeToken(SECRET, NAME, EXPIRY)
	if err != nil {
		t.Fatal(err)
	}

	name, tokenId, expiry, err := readToken(SECRET, token)
	if err != nil {
		t.Fatal(err)
	}

	if name != NAME {
		t.Errorf(`Name mismatch: "%s" != "%s"`, name, NAME)
	}

	if expiry != EXPIRY {
		t.Errorf(`Expiry mismatch: "%d" != "%d"`, expiry, EXPIRY)
	}

	if len(tokenId) == 0 {
		t.Errorf("Expected a token id but it's empty")
	}
}