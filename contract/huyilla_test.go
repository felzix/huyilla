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
    if player.Player.Address == "" {
        t.Error("Expected player to have a private key-associated address but it does not")
    }
}

func TestHuyilla_Login (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    err := h.SignUp(ctx, &types.PlayerName{Name: "felzix"})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    details, err := h.LogIn(ctx, &types.PlayerName{Name: "felzix"})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if details.Player.Spawn == nil {
        t.Error("Player should have gotten a default spawn but did not")
    }

    if details.Entity.Location == nil {
        t.Fatal("Player entity should have been created but it was not (no location)")
    }

    chunk, err := h.GetChunk(ctx, details.Entity.Location.Chunk)
    if err != nil {
        t.Fatalf("Error: %v", err)
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


func TestHuyilla_LoginNegative (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    _, err := h.LogIn(ctx, &types.PlayerName{Name: "felzix"})
    if err == nil {
        t.Fatal("Logging in before signup should throw an error but didn't")
    } else if err.Error() != "Wrong username: no one has this username" {
        t.Errorf("Wrong error. Got %v", err)
    }
}


func TestHuyilla_thisUser (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    addr2 := h.thisUser(ctx)

    if addr1 := addr1.Local.String(); addr1 != addr2 {
        t.Errorf(`Expected addr="%v" but it was "%v"`, addr1, addr2)
    }
}
