package main

import (
    C "github.com/felzix/huyilla/constants"
    "github.com/felzix/huyilla/types"
    "github.com/loomnetwork/go-loom"
    "github.com/loomnetwork/go-loom/plugin"
    "github.com/loomnetwork/go-loom/plugin/contractpb"
    "testing"
)


func TestHuyilla_Chunk (t *testing.T) {
    h := &Huyilla{}

    addr1 := loom.MustParseAddress(ADDR_FROM_LOOM_EXAMPLE)
    ctx := contractpb.WrapPluginContext(plugin.CreateFakeContext(addr1, addr1))

    h.Init(ctx, &plugin.Request{})

    if err := h.GenChunk(ctx, &types.Point{0, 0, 0}); err != nil {
        t.Fatalf(`Error: %v`, err)
    }

    chunk, err := h.GetChunk(ctx, &types.Point{0, 0, 0})
    if err != nil {
        t.Fatalf("Error: %v", err)
    }

    expectedVoxelCount := C.CHUNK_SIZE * C.CHUNK_SIZE * C.CHUNK_SIZE
    if len(chunk.Voxels) != expectedVoxelCount {
        t.Errorf(`Was expected %d voxels but got %d`, expectedVoxelCount, len(chunk.Voxels))
    }
}
