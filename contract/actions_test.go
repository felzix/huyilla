package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)



func TestHuyilla_Actions (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    action := types.Action{
        PlayerName: "admin",
        Action: &types.Action_Move{
            Move: &types.Action_MoveAction{
                WhereTo: &types.AbsolutePoint{
                    &types.Point{1, 12, 144},
                    &types.Point{2, 4, 8},
                },
            },
        },
    }

    // tests behavior when there are no queued actions
    err := h.RegisterAction(ctx, &action)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }
    actions, err := h.getActions(ctx)
    if len(actions.Actions) != 1 {
        t.Errorf(`Expected 1 action but found %d`, len(actions.Actions))
    }

    // tests behavior when there are queued actions
    err = h.RegisterAction(ctx, &action)
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    actions, err = h.getActions(ctx)
    if len(actions.Actions) != 2 {
        t.Errorf(`Expected 2 actions but found %d`, len(actions.Actions))
    }
}
