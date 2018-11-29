package main

import (
	"github.com/felzix/huyilla/types"
	"testing"
)

func TestHuyilla_Content(t *testing.T) {
	h := &Engine{}
	h.Init(&types.Config{})

	content := h.GetContent()

	if content.E[0].Name != "human" {
		t.Errorf(`Expected first entity to be called "human" but it's "%v"`, content.E[0].Name)
	}
}
