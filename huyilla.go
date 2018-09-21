package main

import (
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func main() {
    plugin.Serve(Contract)
}

type Huyilla struct {
}

func (e *Huyilla) Meta() (plugin.Meta, error) {
    return plugin.Meta{
        Name:    "huyilla",
        Version: "0.0.1",
    }, nil
}

func (e *Huyilla) Init(ctx contract.Context, req *plugin.Request) error {
    return nil
}

var Contract plugin.Contract = contract.MakePluginContract(&Huyilla{})
