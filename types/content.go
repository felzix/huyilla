package types

type Content struct {
	E EntityDefinitions
	F FormDefintions
	M MaterialDefinitions
}

type EntityDefinitions map[EntityType]EntityDefinition
type FormDefintions map[Form]FormDefinition
type MaterialDefinitions map[Material]MaterialDefinition

type EntityDefinition struct {
	Name string
	Form
	Material
	Falls bool
	InventoryCapacity uint64
}

type WearableSlot uint8
const (
	WEARABLE_HEAD WearableSlot = iota + 1
	WEARABLE_CHEST
	WEARABLE_HANDS
	WEARABLE_LEGS
	WEARABLE_FEET
	WEARABLE_FINGER
	WEARABLE_NECK
	WEARABLE_EYES
	WEARABLE_WAIST
)

type FormDefinition struct {
	Name string

	Wieldable1Handed bool
	Wieldable2Handed bool
	WearableSlot

	// Tool or Weapon
	//   values are percentages
	Sharpness uint8
	Heft uint8
	Digging uint8
	Mining uint8

	// Clothing or Armor
	//   values are percentages
	PiercingProtection uint8 // with material, for calculating protection from arrows, swords ,etc
	BluntProtection uint8 // with material, for calculating protection from falling, clubs, etc
	Insulation uint8 // with material, for calculating resistance to hot or cold

	// Both
	//   values are percentages
	Strength uint8
}

type MaterialDefinition struct {
	Name string

	// Classifications
	Metallic bool
	Wooden bool
	Fleshy bool
	Salty bool

	// Physical Properties
	Ph uint8
	SolidAt uint32
	GasAt uint32

	// Item Properties
	//   values are percentages
	Sharpness uint8
	Heft uint8
	Strength uint8
}

func (c Content) Marshal() ([]byte, error) {
	return ToBytes(c)
}

func (c *Content) Unmarshal(blob []byte) error {
	return FromBytes(blob, &c)
}
