package engine

type Item interface {
    Form() Form
}


type SimpleItem struct {
    form Form
    material Material
}

type ComplexItem struct {
    form Form
    components []Item
}

type Form uint16
type Material uint16


func MakeSimpleItem (form Form, material Material) *SimpleItem {
    return &SimpleItem{form, material}
}

func MakeComplexItem (form Form, components ...Item) *ComplexItem {
    return &ComplexItem{form, components}
}


func (item *SimpleItem) Form () Form {
    return item.form
}


func (item *SimpleItem) Material () Material {
    return item.material
}


func (item *ComplexItem) Form () Form {
    return item.form
}

func (item *ComplexItem) Components () []Item {
    return item.components
}
