package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)



func TestHuyilla_Entity (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    player, err := h.GetPlayer(ctx, &types.PlayerName{Name: "admin"})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    entity, err := h.GetEntity(ctx, &types.EntityId{Id: player.Player.Id})

    if entity.Type != player.Entity.Type {
        t.Errorf(`GetPlayer and GetEntity returned different entities: "%v" != "%v"`, entity, player.Entity)
    }
}
