package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)



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
