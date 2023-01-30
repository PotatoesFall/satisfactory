package main

type Item string

const (
	IronOre   Item = "Iron Ore"
	CopperOre Item = "Copper Ore"
	Limestone Item = "Limestone"
	Water     Item = "Water"
	Sulfur    Item = "Sulfur"
	Coal      Item = "Coal"

	IronIngot                 Item = "Iron Ingot"
	IronPlate                 Item = "Iron Plate"
	IronRod                   Item = "Iron Rod"
	CopperIngot               Item = "Copper Ingot"
	Wire                      Item = "Wire"
	Cable                     Item = "Cable"
	Concrete                  Item = "Concrete"
	Screw                     Item = "Screw"
	ReinforcedIronPlate       Item = "Reinforced Iron Plate"
	CopperSheet               Item = "Copper Sheet"
	Rotor                     Item = "Rotor"
	ModularFrame              Item = "Modular Frame"
	AssemblyDirectorSystem    Item = "Assembly Director System"
	AdaptiveControlUnit       Item = "Adaptive Control Unit"
	Supercomputer             Item = "Supercomputer"
	Computer                  Item = "Computer"
	HeavyModularFrame         Item = "Heavy Modular Frame"
	CircuitBoard              Item = "Circuit Board"
	AutomatedWiring           Item = "Automated Wiring"
	AILimiter                 Item = "AI Limiter"
	HighSpeedConnector        Item = "High-Speed Connector"
	Plastic                   Item = "Plastic"
	ElectromagneticControlRod Item = "Electromagnetic Control Rod"
	Battery                   Item = "Battery"
	Stator                    Item = "Stator"
	EncasedIndustrialBeam     Item = "Encased Industrial Beam"
	SteelPipe                 Item = "Steel Pipe"
	AlcladAluminumSheet       Item = "Alclad Aluminum Sheet"
	Quickwire                 Item = "Quickwire"
	SteelIngot                Item = "Steel Ingot"
	CateriumIngot             Item = "Caterium Ingot"
	CrystalOscillator         Item = "Crystal Oscillator"
	CateriumOre               Item = "Caterium Ore"
	Silica                    Item = "Silica"
	RawQuartz                 Item = "Raw Quartz"
	QuartzCrystal             Item = "Quartz Crystal"
	Rubber                    Item = "Rubber"
	AluminumIngot             Item = "Aluminum Ingot"
	MagneticFieldGenerator    Item = "Magnetic Field Generator"
	VersatileFramework        Item = "Versatile Framework"
	SteelBeam                 Item = "Steel Beam"
	PolymerResin              Item = "Polymer Resin"
	CrudeOil                  Item = "Crude Oil"
	HeavyOilResidue           Item = "Heavy Oil Residue"
)

var rawItems = [...]Item{
	IronOre,
	CopperOre,
	Limestone,
	Water,
	Sulfur,
	Coal,
	CateriumOre,
	RawQuartz,
	CrudeOil,
}

func isRaw(i Item) bool {
	for _, raw := range rawItems {
		if i == raw {
			return true
		}
	}

	return false
}

func isBase(i Item) bool {
	return i == AluminumIngot || isRaw(i)
}
