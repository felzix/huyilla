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

type MaterialProperty struct {
    // Classifications
    Metallic bool
    Wooden   bool
    Fleshy   bool
    Salty    bool

    // Item Properties
    Sharpness Percentage  // for calculating piercing damage
    Heft      Percentage  // for calculating blunt damage; not used for mass
    Strength  Percentage  // for calculating how much damage the item can sustain

    // Physical Properties
    PH      uint8
    SolidAt Kelvin
    GasAt   Kelvin
}
type FormProperty struct {
    // How it's used
    Wieldable1Handed bool
    Wieldable2Handed bool

    WearableHead   bool
    WearableChest  bool
    WearableHands  bool
    WearableLegs   bool
    WearableFeet   bool
    WearableFinger bool
    WearableNeck   bool
    WearableEyes   bool
    WearableWaist  bool

    // Tool or Weapon
    Sharpness Percentage
    Heft      Percentage
    Digging   Percentage
    Mining    Percentage

    // Clothing or Armor
    PiercingProtection Percentage  // with material, for calculating protection from arrows, swords ,etc
    BluntProtection    Percentage  // with material, for calculating protection from falling, clubs, etc
    Insulation         Percentage  // with material, for calculating resistance to hot or cold

    // Both
    Strength Percentage
}

type Kelvin uint16
type Percentage uint8


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
