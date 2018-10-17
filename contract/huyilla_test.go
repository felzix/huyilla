package contract

import (
    "testing"
)

const ADDR_FROM_LOOM_EXAMPLE = "chain:0xb16a379ec18d4093666f8f38b11a3071c920207d"

func Test_Huyilla_Meta (t *testing.T) {
    h := &Huyilla{}
    meta, err := h.Meta()
    if err != nil {
        t.Fatalf(`Error: %v`, err)
    }

    if meta.Name != "Huyilla" {
        t.Errorf(`Contract name is "%v"; should be "Huyilla"`, meta.Name)
    }
    if meta.Version != "0.0.1" {
        t.Errorf(`Contract version is "%v"; should be "0.0.1"`, meta.Version)
    }
}
