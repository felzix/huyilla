package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)

const ADDR_FROM_LOOM_EXAMPLE = "chain:0xb16a379ec18d4093666f8f38b11a3071c920207d"

func Test_Huyilla_Meta (t *testing.T) {
    h := &Huyilla{}
    meta, err := h.Meta()
    if err != nil {
        t.Fatalf(`Error: %v`, err)
    }

    if meta.Name != "Huyilla" {
        t.Errorf(`Contract name is "%v"; should be "Huyilla"`, meta.Name)
    }
    if meta.Version != "0.0.1" {
        t.Errorf(`Contract version is "%v"; should be "0.0.1"`, meta.Version)
    }
}


func TestHuyilla_SignUp (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    err := h.SignUp(ctx, &types.PlayerName{Name: "felzix"})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    player, err := h.GetPlayer(ctx, &types.PlayerName{Name: "felzix"})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if player.Player.Name != player.Entity.PlayerName {
        t.Errorf(
            `Player's entity has wrong name: player="%v", entity="%v"`,
            player.Player.Name, player.Entity.PlayerName)
    }
    if player.Player.Address == nil {
        t.Error("Expected player to have a private key-associated address but it does not")
    }
}
