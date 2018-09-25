package engine

type MapGenerator interface {
    GenerateChunk(Point)
}

