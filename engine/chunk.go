package engine


type Chunk struct {
	voxels   VoxelCube
	entities Entities
	items    Items
}
type VoxelCube [CHUNK_SIZE][CHUNK_SIZE][CHUNK_SIZE]Voxel
type Entities map[UniqueId]EntityAt
type EntityAt struct{
    Who   *Entity
    Where At
}
type Items []ItemAt
type ItemAt struct {
    Which *Item
    Where At
}
type At struct{
    X float32
    Y float32
    Z float32
}

func MakeChunk (voxels VoxelCube) *Chunk {
    return &Chunk{voxels, make(Entities, 0), make(Items, 0)}
}

func (chunk *Chunk) Get (p Point) (*Voxel) {
    return &chunk.voxels[p.Y][p.X][p.Z]
}


func (chunk *Chunk) Set (p Point, v Voxel) {
    chunk.voxels[p.Y][p.X][p.Z] = v
}

func (chunk *Chunk) GetEntity (who UniqueId) *Entity {
    return chunk.entities[who].Who
}

func (chunk *Chunk) SetEntity (entity *Entity, where At) {
    chunk.entities[entity.Id] = EntityAt{Who: entity, Where: where}
}

func (chunk *Chunk) RemoveEntity (who UniqueId) {
    delete(chunk.entities, who)
}

func (chunk *Chunk) GetItems () Items {
    return chunk.items
}

func (chunk *Chunk) AddItem (item *Item, where At) {
    chunk.items = append(chunk.items, ItemAt{item, where})
}
