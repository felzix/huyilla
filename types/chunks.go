package types

type Chunks struct {
	Chunks map[Point]Chunk
	Entities map[EntityId]Entity
	Items map[ItemId]Item
}

func NewChunks(radius uint64) *Chunks {
	diameter := 1 + radius*2
	size := diameter * diameter * diameter
	return &Chunks{
		Chunks: make(map[Point]Chunk, size),
		Entities: make(map[EntityId]Entity, 0),
		Items: make(map[ItemId]Item, 0),
	}
}

func (c Chunks) Marshal() ([]byte, error) {
	return ToBytes(c)
}

func (c *Chunks) Unmarshal(blob []byte) error {
	return FromBytes(blob, &c)
}
