package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)



func TestHuyilla_ActiveChunkRadius (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    err := h.SignUp(ctx, &types.PlayerName{"felzix"})
    if err != nil {
        t.Fatal(err)
    }
    player, err := h.LogIn(ctx, &plugin.Request{})
    if err != nil {
        t.Fatal(err)
    }
    err = h.Tick(ctx, &plugin.Request{})
    if err != nil {
        t.Fatal(err)
    }

    chunk, err := h.getChunk(ctx, player.Entity.Location.Chunk)
    if err != nil {
        t.Fatal(err)
    }
    if len(chunk.Entities) != 1 {
        t.Errorf("Expected 1 entity in chunk but there were %d", len(chunk.Entities))
    }

    // active range in the positive direction
    edge := clonePoint(player.Entity.Location.Chunk)
    edge.X += 3
    chunk, err = h.getChunk(ctx, edge)
    if chunk == nil {
        t.Error("Chunk within player's range should exist but it does not.")
    }
    if len(chunk.Entities) != 0 {
        t.Errorf("Expected 0 entities in chunk but there were %d", len(chunk.Entities))
    }

    beyond := clonePoint(player.Entity.Location.Chunk)
    beyond.X += 4
    chunk, err = h.getChunk(ctx, beyond)
    if err == nil {  // note that this is "==" not "!="
        t.Error("Chunk beyond player's range exists when it should not.")
    }

    // active range in the negative direction
    edge = clonePoint(player.Entity.Location.Chunk)
    edge.X -= 3
    chunk, err = h.getChunk(ctx, edge)
    if err != nil {
        t.Fatal(err)
    }
    if chunk == nil {
        t.Error("Chunk within player's range should exist but it does not.")
    }

    beyond = clonePoint(player.Entity.Location.Chunk)
    beyond.X -= 4
    chunk, err = h.getChunk(ctx, beyond)
    if err == nil {  // note that this is "==" not "!="
        t.Error("Chunk beyond player's range exists when it should not.")
    }
}
