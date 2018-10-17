package main

import (
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)

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

    addr1 := loom.MustParseAddress("chain:0xb16a379ec18d4093666f8f38b11a3071c920207d")
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
