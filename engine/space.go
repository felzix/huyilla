package engine


func NewSparseSpace() *SparseSpace {
    return &SparseSpace{
        space:  make(map[int]map[int]map[int]interface{}),
        length: 0,
    }
}


type SparseSpace struct {
    space  map[int]map[int]map[int]interface{}
    length uint
}

func (s *SparseSpace) Set (p Point, element interface{}) {
    if _, ok := s.space[p.X]; !ok {
        s.space[p.X] = make(map[int]map[int]interface{})
    }
    if _, ok := s.space[p.X][p.Y]; !ok {
        s.space[p.X][p.Y] = make(map[int]interface{})
    }
    if _, ok := s.space[p.X][p.Y][p.Z]; !ok {
        s.length++  // new element
    }

    s.space[p.X][p.Y][p.Z] = element
}

func (s *SparseSpace) Get (p Point) interface{} {
    return s.space[p.X][p.Y][p.Z]
}

func (s *SparseSpace) Delete (p Point) {
    delete(s.space[p.X][p.Y], p.Z)

    // one fewer element
    s.length--

    // clean up
    if line := s.space[p.X][p.Y]; len(line) == 0 {
        delete(s.space[p.X], p.Y)

        if plane := s.space[p.X]; len(plane) == 0 {
            delete(s.space, p.X)

        }
    }
}