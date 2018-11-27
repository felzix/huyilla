package main

import (
	"github.com/felzix/huyilla/types"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/go-loom/plugin/contractpb"
	"testing"
)

func TestHuyilla_Players(t *testing.T) {
	h := &Huyilla{}

	addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
	ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

	h.Init(ctx, &plugin.Request{})

	h.SignUp(ctx, &types.PlayerName{"felzix"})
	h.LogIn(ctx, &plugin.Request{})

	players, err := h.GetPlayerList(ctx, &plugin.Request{})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(players.Names) != 2 { // FAKE and new player "felzix"
		t.Errorf(`Error: Should be two players but there aren't: "%v"`, players.Names)
	}

	player, err := h.GetPlayer(ctx, &types.Address{addr1.Local.String()})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.Player.Name != "felzix" {
		t.Errorf(`Player name was "%v" instead of "felzix"`, player.Player.Name)
	}
}
