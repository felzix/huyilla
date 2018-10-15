package main

import "testing"

func TestHuyilla_GetConfig(t *testing.T) {
    meta, err := Contract.Meta()
    if err != nil {
        t.Errorf(`Error: %v`, err)
    }

    if meta.Name != "Huyilla" {
        t.Errorf(`Contract name is "%v"; should be "Huyilla"`, meta.Name)
    }
    if meta.Version != "0.0.1" {
        t.Errorf(`Contract version is "%v"; should be "0.0.1"`, meta.Version)
    }
}