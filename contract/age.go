package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)


var AGE = []byte("Age")


func (c *Huyilla) GetAge (ctx contract.StaticContext, req *plugin.Request) (*types.Age, error) {
    var age = &types.Age{}
    err := ctx.Get(AGE, age)
    return age, err
}

func (c *Huyilla) incrementAge (ctx contract.Context, req *plugin.Request) (*types.Age, error) {
    age, err := c.GetAge(ctx, req)
    if err != nil { return nil, err }

    age.Ticks ++
    err = ctx.Set(AGE, age)
    return age, err
}
