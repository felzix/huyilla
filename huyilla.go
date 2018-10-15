package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


func main() {
    plugin.Serve(Contract)
}

type Huyilla struct {}


var Contract = contract.MakePluginContract(&Huyilla{})


func (c *Huyilla) Meta () (plugin.Meta, error) {
    return plugin.Meta{
        Name:    "Huyilla",
        Version: "0.0.1",
    }, nil
}

var AGE = []byte("Age")
var CONFIG = []byte("Config")
var PLAYERS = []byte("Players")

func (c *Huyilla) Init (ctx contract.Context, req *plugin.Request) error {
    ctx.Set(AGE, &types.Age{Ticks: 1})  // starts at 1 because 0 counts as non-existent
    return nil
}

func (c *Huyilla) GetAge (ctx contract.StaticContext, req *plugin.Request) (*types.Age, error) {
    var age = &types.Age{}
    err := ctx.Get(AGE, age)
    return age, err
}

func (c *Huyilla) incrementAge (ctx contract.Context, req *plugin.Request) (*types.Age, error) {
    age, err := c.GetAge(ctx, req)
    if err != nil { return age, err }

    age.Ticks ++
    err = ctx.Set(AGE, age)
    return age, err
}

// func (c *Huyilla) SetAge (ctx contract.Context, req *types.Age) error {
//     return ctx.Set(AGE, req)
// }