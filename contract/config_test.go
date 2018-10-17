package contract

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)


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
