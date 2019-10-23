package main

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"sync"
)

type Cache struct {
	sync.Mutex

	scene          *core.Node
	age            uint64
	previousChunks map[types.ComparablePoint]*types.DetailedChunk
	chunks         map[types.ComparablePoint]*types.DetailedChunk
	meshes         map[types.ComparablePoint]*Meshes
}

type Meshes struct {
	Voxels   []*graphic.Mesh
	Entities map[int64]*graphic.Mesh
}

func (meshes *Meshes) GetVoxelIndex(x, y, z int) int {
	return types.CalculateVoxelIndex(uint64(x), uint64(y), uint64(z))
}

func (meshes *Meshes) GetVoxel(x, y, z int) *graphic.Mesh {
	i := meshes.GetVoxelIndex(x, y, z)
	return meshes.Voxels[i]
}

func (meshes *Meshes) SetVoxel(x, y, z int, mesh *graphic.Mesh) {
	i := meshes.GetVoxelIndex(x, y, z)
	oldMesh := meshes.Voxels[i]
	if oldMesh != nil {
		mesh.Dispose()
	}
	meshes.Voxels[i] = mesh
}

func (cache *Cache) DestroyEntity(id int64, point types.ComparablePoint, guiClient *GuiClient) {
	entity := cache.chunks[point].Entities[id]
	if entity.PlayerName == guiClient.player.Player.Name {
		guiClient.playerNode = nil
		// guiClient.camera // TODO do something with the camera
	}

	meshes := cache.meshes[point]
	meshes.Entities[id].Dispose()
	delete(meshes.Entities, id)
}

func (cache *Cache) MoveEntity(
	id int64,
	point types.ComparablePoint,
	offset *types.Point) {
	entity := cache.chunks[point].Entities[id]
	mesh := cache.meshes[point].Entities[id]
	x := float32(entity.Location.Voxel.X + (offset.X * 16))
	y := float32(entity.Location.Voxel.Y + (offset.Y * 16))
	z := float32(entity.Location.Voxel.Z + (offset.Z * 16))
	mesh.SetPosition(x, y, z)
}

func (cache *Cache) DrawEntity(
	entity *types.Entity,
	point types.ComparablePoint,
	offset *types.Point,
	guiClient *GuiClient) {
	def := content.EntityDefinitions[entity.Type]
	geom := geometries[def.Form]
	mat := materials[def.Material]

	mesh := graphic.NewMesh(geom, mat)
	x := float32(entity.Location.Voxel.X + (offset.X * 16))
	y := float32(entity.Location.Voxel.Y + (offset.Y * 16))
	z := float32(entity.Location.Voxel.Z + (offset.Z * 16))
	mesh.SetPosition(x, y, z)
	cache.scene.Add(mesh)
	mesh.SetRotation(math32.Pi/2, 0, 0)
	cache.meshes[point].Entities[entity.Id] = mesh

	if entity.PlayerName == guiClient.player.Player.Name {
		fmt.Println("create new player mesh")
		guiClient.playerNode = mesh
		mesh.Add(guiClient.camera)
	}
}

func (meshes *Meshes) ClearVoxel(x, y, z int) {
	i := meshes.GetVoxelIndex(x, y, z)
	meshes.Voxels[i] = nil
}

func NewCache(scene *core.Node) *Cache {
	return &Cache{
		scene:          scene,
		age:            0,
		chunks:         make(map[types.ComparablePoint]*types.DetailedChunk, 0),
		previousChunks: make(map[types.ComparablePoint]*types.DetailedChunk, 0),
		meshes:         make(map[types.ComparablePoint]*Meshes, 0),
	}
}

func (cache *Cache) GetAge() uint64 {
	return cache.age
}

func (cache *Cache) SetAge(age uint64) {
	cache.Lock()
	defer cache.Unlock()
	cache.age = age
}

func (cache *Cache) AddMeshes(point types.ComparablePoint) *Meshes {
	meshes := &Meshes{
		Voxels: make([]*graphic.Mesh, C.CHUNK_LENGTH),
		Entities: make(map[int64]*graphic.Mesh, 0),
	}
	cache.meshes[point] = meshes
	return meshes
}

func (cache *Cache) DeleteMeshes(point types.ComparablePoint) {
	mesh := cache.meshes[point]
	for _, voxel := range mesh.Voxels {
		if voxel != nil {
			voxel.Dispose()
		}
	}
	for _, entity := range mesh.Entities {
		entity.Dispose()
	}
	delete(cache.meshes, point)
}

func (cache *Cache) GetChunk(coords *types.Point) *types.DetailedChunk {
	point := *types.NewComparablePoint(coords)
	return cache.chunks[point]
}

func (cache *Cache) GetPreviousChunk(coords *types.Point) *types.DetailedChunk {
	point := *types.NewComparablePoint(coords)
	return cache.previousChunks[point]
}

func (cache *Cache) SetChunk(coords *types.Point, chunk *types.DetailedChunk) {
	cache.Lock()
	defer cache.Unlock()
	point := *types.NewComparablePoint(coords)
	if cache.chunks[point] != nil {
		cache.previousChunks[point] = cache.chunks[point]
	}

	cache.chunks[point] = chunk
}

func EachVoxel(fn func(int, int, int)) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				fn(x, y, z)
			}
		}
	}
}

func (cache *Cache) DrawVoxel(
	voxel types.Voxel,
	point types.ComparablePoint,
	x, y, z int,
	offset *types.Point,
	meshes *Meshes,
	) {
	if isDrawn(voxel) {
		v := voxel.Expand()
		geom := geometries[v.Form]
		mat := materials[v.Material]
		mesh := graphic.NewMesh(geom, mat)
		meshes.SetVoxel(x, y, z, mesh)
		cache.scene.Add(mesh)
		sceneX := float32(x + int(offset.X*16))
		sceneY := float32(y + int(offset.Y*16))
		sceneZ := float32(z + int(offset.Z*16))
		mesh.SetPosition(sceneX, sceneY, sceneZ)
	} else {
		meshes.ClearVoxel(x, y, z)
	}
}

func (cache *Cache) Draw(guiClient *GuiClient) {
	center := guiClient.player.Entity.Location.Chunk

	// Creates, alters or removes meshes as needed.

	for point, previousChunk := range cache.previousChunks {
		currentChunk := cache.chunks[point]
		meshes := cache.meshes[point]
		if currentChunk == nil {
			cache.DeleteMeshes(point)
		} else {
			offset := center.DeriveVector(point.ToPoint())
			EachVoxel(func(x, y, z int) {
				previousVoxel := previousChunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				voxel := currentChunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				if voxel != previousVoxel {
					cache.DrawVoxel(voxel, point, x, y, z, offset, meshes)
				}
			})

			for id, previousEntity := range previousChunk.Entities {
				currentEntity := currentChunk.Entities[id]
				if currentEntity == nil {
					cache.DestroyEntity(id, point, guiClient)
				} else {
					if !previousEntity.Location.Equals(currentEntity.Location) {
						cache.MoveEntity(id, point, offset)
					}
				}
			}
			for id, entity := range currentChunk.Entities {
				previousEntity := previousChunk.Entities[id]
				if previousEntity == nil {
					cache.DrawEntity(entity, point, offset, guiClient)
				}
			}
		}
	}
	for point, currentChunk := range cache.chunks {
		previousChunk := cache.previousChunks[point]
		if previousChunk == nil {
			offset := center.DeriveVector(point.ToPoint())
			meshes := cache.AddMeshes(point)
			EachVoxel(func(x, y, z int) {
				voxel := currentChunk.GetVoxel(uint64(x), uint64(y), uint64(z))
				cache.DrawVoxel(voxel, point, x, y, z, offset, meshes)
			})
			for _, entity := range currentChunk.Entities {
				cache.DrawEntity(entity, point, offset, guiClient)
			}
		}
	}
}
