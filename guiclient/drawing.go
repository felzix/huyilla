package main

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

func (guiClient *GuiClient) buildVoxels(previousChunk, chunk *types.DetailedChunk, offset *types.Point) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				previousVoxel := types.Voxel(0)
				if previousChunk == nil {
					previousVoxel = 0
				} else {
					previousVoxel = previousChunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				}

				voxel := chunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				if voxel != previousVoxel && isDrawn(voxel) {
					trueX := float32(x + int(offset.X*16))
					trueY := float32(y + int(offset.Y*16))
					trueZ := float32(z + int(offset.Z*16))
					makeVoxel(guiClient.rootScene, trueX, trueY, trueZ, voxel)
				}
			}
		}
	}
	for _, e := range chunk.Entities {
		eX := float32(e.Location.Voxel.X + (offset.X * 16))
		eY := float32(e.Location.Voxel.Y + (offset.Y * 16))
		eZ := float32(e.Location.Voxel.Z + (offset.Z * 16))
		guiClient.makeEntity(eX, eY, eZ, e)
	}
}

func (guiClient *GuiClient) makeEntity(x, y, z float32, entity *types.Entity) {
	def := content.EntityDefinitions[entity.Type]
	geom := geometries[def.Form]
	mat := materials[def.Material]

	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x, y, z+1)
	guiClient.rootScene.Add(mesh)
	mesh.SetRotation(math32.Pi/2, 0, 0)

	if entity.PlayerName == guiClient.player.Player.Name {
		guiClient.playerNode = mesh
		mesh.Add(guiClient.camera)
	}
}

func isDrawn(voxel types.Voxel) bool {
	M := content.MATERIAL
	v := voxel.Expand()

	valid := map[uint64]bool {
		M["dirt"]: true,
		M["grass"]: true,
		M["water"]: true,
	}

	return valid[v.Material]
}

func makeVoxel(scene *core.Node, x, y, z float32, voxel types.Voxel) {
	v := voxel.Expand()
	geom := geometries[v.Form]
	mat := materials[v.Material]

	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)
	mesh.SetPosition(x, y, z)
}
