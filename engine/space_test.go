package engine

import (
    "testing"
)

func TestSparseSpace(t *testing.T) {
    s := NewSparseSpace()

    s.Set(Point{0, 0, 0}, 19)
    if v := s.Get(Point{0, 0, 0}); v != 19 {
        t.Errorf(`SparseSpace stored "%v" instead of "%v"`, v, 19)
    }
}

