package main

import (
	C "github.com/felzix/huyilla/constants"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"sync"
)

type Cache struct {
	sync.Mutex

	age              types.Age
	previousChunks   map[types.Point]*types.Chunk
	chunks           map[types.Point]*types.Chunk
	previousEntities map[types.EntityId]*types.Entity

	scene        *core.Node
	entityMeshes map[types.EntityId]*graphic.Mesh
	voxelMeshes  map[types.Point][]*graphic.Mesh // pt -> [mesh]
	basis        *types.Point
}

func NewCache(scene *core.Node) *Cache {
	return &Cache{
		age:              types.Age(0),
		previousChunks:   make(map[types.Point]*types.Chunk, 0),
		chunks:           make(map[types.Point]*types.Chunk, 0),
		previousEntities: make(map[types.EntityId]*types.Entity),

		scene:        scene,
		voxelMeshes:  make(map[types.Point][]*graphic.Mesh),
		entityMeshes: make(map[types.EntityId]*graphic.Mesh),
		basis:        nil,
	}
}

func (cache *Cache) GetAge() types.Age {
	return cache.age
}

func (cache *Cache) SetAge(age types.Age) {
	cache.Lock()
	defer cache.Unlock()
	cache.age = age
}

func (cache *Cache) SetChunks(chunks *types.Chunks) {
	// TODO do the same thing here as you do with meshes
	// Possibly combine this and Draw(), or rearrange them, so the client just gives the cache an update
	// No need, I think, for the client to call SetChunks then call Draw.
	for point, chunk := range chunks.Chunks {
		cache.setChunk(point, &chunk)
	}
}

func (cache *Cache) setChunk(point types.Point, chunk *types.Chunk) {
	if cache.chunks[point] != nil {
		cache.previousChunks[point] = cache.chunks[point]
	}

	cache.chunks[point] = chunk
}

func (cache *Cache) CreateVoxelMeshes(point types.Point, chunk *types.Chunk, offset types.Point) {
	Log.Debug("CreateVoxelMeshes", point.ToString())
	cache.voxelMeshes[point] = make([]*graphic.Mesh, C.CHUNK_LENGTH)
	types.EachVoxel(func(x, y, z uint64) {
		voxel := chunk.GetVoxel(x, y, z)
		cache.DrawVoxel(voxel, point, x, y, z, offset)
	})
}

func (cache *Cache) DestroyVoxelMeshes(point types.Point) {
	Log.Debug("DestroyVoxelMeshes", point.ToString())
	types.EachVoxel(func(x, y, z uint64) {
		i := types.CalculateVoxelIndex(x, y, z)
		mesh := cache.voxelMeshes[point][i]
		if mesh != nil {
			mesh.Dispose()
		}
	})
	delete(cache.voxelMeshes, point)
}

func (cache *Cache) UpdateVoxelMeshes(
	point types.Point,
	offset types.Point,
	previousChunk *types.Chunk,
	currentChunk *types.Chunk,
) {
	types.EachVoxel(func(x, y, z uint64) {
		previousVoxel := previousChunk.GetVoxel(x, y, z)
		voxel := currentChunk.GetVoxel(x, y, z)
		if voxel != previousVoxel {
			Log.Debug("UpdateVoxelMeshes", point.ToString(), x, y, z)
			cache.DrawVoxel(voxel, point, x, y, z, offset)
		}
	})
}

func (cache *Cache) CreateEntityMesh(entity *types.Entity, offset types.Point, guiClient *GuiClient) {
	Log.Debug("CreateEntityMesh", entity.Id)
	def := content.EntityDefinitions[entity.Type]
	geom := geometries[def.Form]
	mat := materials[def.Material]
	mesh := graphic.NewMesh(geom, mat)
	cache.entityMeshes[entity.Id] = mesh

	SetMeshPosition(mesh, entity.Location.Voxel, offset)

	cache.scene.Add(mesh)

	if entity.PlayerName == guiClient.player.Player.Name {
		guiClient.playerNode = mesh
		mesh.Add(guiClient.camera)
		// Sets camera a bit higher. Unnecessary once player meshes have heads.
		guiClient.camera.SetPosition(0, 1, 0)
	}
}

func (cache *Cache) DestroyEntityMesh(entity *types.Entity, guiClient *GuiClient) {
	Log.Debug("DestroyEntityMesh", entity.Id)
	if entity.PlayerName == guiClient.player.Player.Name && cache.entityMeshes[entity.Id] == guiClient.playerNode {
		guiClient.playerNode = nil
		// TODO do something with camera
	}

	cache.entityMeshes[entity.Id].Dispose()
	delete(cache.entityMeshes, entity.Id)
}

func (cache *Cache) UpdateEntityMesh(previousEntity *types.Entity, entity *types.Entity, offset types.Point) {
	if previousEntity.Location.Equals(entity.Location) {
		return
	}
	Log.Debug("UpdateEntityMesh", entity.Id)
	mesh := cache.entityMeshes[entity.Id]
	SetMeshPosition(mesh, entity.Location.Voxel, offset)
}

func (cache *Cache) DrawVoxel(
	voxel types.Voxel,
	point types.Point,
	x, y, z uint64,
	offset types.Point,
) {
	v := voxel.Expand()
	if isDrawn(v) {
		geom := geometries[v.Form]
		mat := materials[v.Material]
		mesh := graphic.NewMesh(geom, mat)
		cache.SetVoxelMesh(point, x, y, z, mesh)

		cache.scene.Add(mesh)
		SetMeshPosition(mesh, types.NewPoint(int64(x), int64(y), int64(z)), offset)
	} else {
		cache.SetVoxelMesh(point, x, y, z, nil)
	}
}

func (cache *Cache) SetVoxelMesh(point types.Point, x, y, z uint64, mesh *graphic.Mesh) {
	i := types.CalculateVoxelIndex(x, y, z)
	oldMesh := cache.voxelMeshes[point][i]
	if oldMesh != nil {
		oldMesh.Dispose()
	}
	cache.voxelMeshes[point][i] = mesh
}

// Creates, alters or removes meshes as needed.
func (cache *Cache) Draw(guiClient *GuiClient) {
	cache.Lock()
	defer cache.Unlock()

	// TODO handle a moving basis
	if cache.basis == nil {
		cache.basis = &guiClient.player.Entity.Location.Chunk
	}
	center := *cache.basis
	currentEntities := make(map[types.EntityId]*types.Entity, 0)

	for point, currentChunk := range cache.chunks {
		previousChunk := cache.previousChunks[point]
		if previousChunk == nil {
			offset := Offset(center, point)
			cache.CreateVoxelMeshes(point, currentChunk, offset)
			for _, id := range currentChunk.Entities {
				currentEntities[id] = cache.previousEntities[id]
			}
		}
	}
	for point, previousChunk := range cache.previousChunks {
		currentChunk := cache.chunks[point]
		if currentChunk == nil {
			cache.DestroyVoxelMeshes(point)
		} else {
			offset := Offset(center, point)
			cache.UpdateVoxelMeshes(point, offset, previousChunk, currentChunk)
			for _, id := range currentChunk.Entities {
				currentEntities[id] = cache.previousEntities[id]
			}
		}
	}

	for id, currentEntity := range currentEntities {
		previousEntity := cache.previousEntities[id]
		if previousEntity == nil {
			offset := Offset(center, currentEntity.Location.Chunk)
			cache.CreateEntityMesh(currentEntity, offset, guiClient)
		}
	}
	for id, previousEntity := range cache.previousEntities {
		currentEntity := currentEntities[id]
		if currentEntity == nil {
			cache.DestroyEntityMesh(previousEntity, guiClient)
		} else {
			offset := Offset(center, currentEntity.Location.Chunk)
			cache.UpdateEntityMesh(previousEntity, currentEntity, offset)
		}
	}

	cache.previousEntities = currentEntities
}

func Offset(basis, point types.Point) types.Point {
	return basis.DeriveVector(point)
}

func SetMeshPosition(mesh *graphic.Mesh, point types.Point, offset types.Point) {
	// mixing y and z is intended here
	x := float32(point.X + (offset.X * 16))
	y := float32(point.Z + (offset.Z * 16))
	z := float32(point.Y + (offset.Y * 16))
	mesh.SetPosition(x, y, z)
}

func isDrawn(v types.ExpandedVoxel) bool {
	M := content.MATERIAL

	valid := map[types.Material]bool{
		M["dirt"]:  true,
		M["grass"]: true,
		M["water"]: true,
	}

	return valid[v.Material]
}
