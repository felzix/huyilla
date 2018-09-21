package engine


func NewSparseSpace() *SparseSpace {
    return &SparseSpace{
        space: make(map[int]map[int]map[int]interface{}),
        length: 0,
    }
}


type SparseSpace struct {
    space map[int]map[int]map[int]interface{}
    length uint
}

func (s *SparseSpace) Set (p Point, element interface{}) {
    if _, ok := s.space[p.x]; !ok {
        s.space[p.x] = make(map[int]map[int]interface{})
    }
    if _, ok := s.space[p.x][p.y]; !ok {
        s.space[p.x][p.y] = make(map[int]interface{})
    }
    if _, ok := s.space[p.x][p.y][p.z]; !ok {
        s.length++  // new element
    }

    s.space[p.x][p.y][p.z] = element
}

func (s *SparseSpace) Get (p Point) interface{} {
    return s.space[p.x][p.y][p.z]
}

func (s *SparseSpace) Delete (p Point) {
    delete(s.space[p.x][p.y], p.z)

    // one fewer element
    s.length--

    // clean up
    if line := s.space[p.x][p.y]; len(line) == 0 {
        delete(s.space[p.x], p.y)

        if plane := s.space[p.x]; len(plane) == 0 {
            delete(s.space, p.x)

        }
    }
}