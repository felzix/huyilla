package contract

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

func TestHuyilla_Age (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    age, err := h.GetAge(ctx, &plugin.Request{})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    if age.Ticks != 1 {
        t.Errorf(`Expected age to be the default of "1" not "%v"`, age.Ticks)
    }
}

func TestHuyilla_Config (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    h.SetConfigOptions(ctx, &types.PrimitiveMap{
        Map: map[string]*types.Primitive{
            "PlayerCap": {Value: &types.Primitive_Int{Int: 101}},
        }})

    config, err := h.GetConfig(ctx, &plugin.Request{})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    foundPlayerCap := config.Options.Map["PlayerCap"].GetInt()
    if foundPlayerCap != 101 {
        t.Errorf(`Expected player cap to be "101" not "%v"`, foundPlayerCap)
    }
}


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
