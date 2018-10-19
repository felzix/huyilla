package main

import (
    "github.com/felzix/huyilla/content"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom/plugin"
    contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)



func (c *Huyilla) GetContent (ctx contract.StaticContext, req *plugin.Request) (*types.Content, error) {
    return &content.Content, nil
}
