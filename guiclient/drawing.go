package main
import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

func buildVoxels(scene *core.Node, chunk *types.DetailedChunk, offset *types.Point) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				voxel := chunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				if isDrawn(voxel) {
					trueX := float32(x + int(offset.X * 16))
					trueY := float32(y + int(offset.Y * 16))
					trueZ := float32(z + int(offset.Z * 16))
					makeVoxel(scene, trueX, trueY, trueZ, voxel)
				}
			}
		}
	}
	for _, e := range chunk.Entities {
		eX := float32(e.Location.Voxel.X + (offset.X * 16))
		eY := float32(e.Location.Voxel.Y + (offset.Y * 16))
		eZ := float32(e.Location.Voxel.Z + (offset.Z * 16))
		makeEntity(scene, eX, eY, eZ, e)
	}
}

func makeEntity(scene *core.Node, x, y, z float32, entity *types.Entity) {
	def := content.EntityDefinitions[entity.Type]
	geom := geometries[def.Form]
	mat := materials[def.Material]

	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(x, y, z + 1)
	scene.Add(mesh)
	mesh.SetRotation(math32.Pi/2, 0, 0)
}

func isDrawn(voxel types.Voxel) bool {
	M := content.MATERIAL
	v := voxel.Expand()
	switch v.Material {
	case M["air"]:
		return false
	case M["dirt"]:
		return true
	case M["grass"]:
		return true
	case M["water"]:
		return true
	default:
		return false
	}
}

func makeVoxel(scene *core.Node, x, y, z float32, voxel types.Voxel) {
	v := voxel.Expand()
	geom := geometries[v.Form]
	mat := materials[v.Material]

	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)
	mesh.SetPosition(x, y, z)
}
