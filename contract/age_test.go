package contract

import (
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)


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
