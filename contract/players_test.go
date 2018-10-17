package contract

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)



func TestHuyilla_Player (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    players, err := h.GetPlayerList(ctx, &plugin.Request{})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if len(players.Names) != 1 {  // default has admin
        t.Errorf(`Error: Should be no players but there aren't: "%v"`, players.Names)
    }

    player, err := h.GetPlayer(ctx, &types.PlayerName{Name: "admin"})

    if err != nil && err.Error() != "not found" {
        t.Errorf(`Expected "not found" error but got "%v"`, err)
    } else if err == nil && player != nil {
        t.Errorf(`Expected "not found" error but player was actually returned: "%v"`, player)
    } else if err == nil && player == nil {
        t.Error(`Expected "not found" error but instead nil entity was returned`)
    }
}
