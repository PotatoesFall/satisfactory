package main

type Recipe struct {
	Name   string
	Input  map[Item]uint
	Output map[Item]uint
	// Energy uint
}

var allRecipes = map[Item]Recipe{
	IronIngot: {
		Name: "Pure Iron Ingot",
		Input: map[Item]uint{
			IronOre: 7,
			Water:   4,
		},
		Output: map[Item]uint{
			IronIngot: 13,
		},
	},
	AssemblyDirectorSystem: {
		Name: "Assembly Director System",
		Input: map[Item]uint{
			AdaptiveControlUnit: 2,
			Supercomputer:       1,
		},
		Output: map[Item]uint{
			AssemblyDirectorSystem: 1,
		},
	},
	AdaptiveControlUnit: {
		Name: "Adaptive Control Unit",
		Input: map[Item]uint{
			AutomatedWiring:   15,
			CircuitBoard:      10,
			HeavyModularFrame: 2,
			Computer:          2,
		},
		Output: map[Item]uint{
			AdaptiveControlUnit: 2,
		},
	},
	Supercomputer: {
		// Name: "Supercomputer",
		// Input: map[Item]uint{
		// 	Computer:           2,
		// 	AILimiter:          2,
		// 	HighSpeedConnector: 3,
		// 	Plastic:            28,
		// },
		// Output: map[Item]uint{
		// 	Supercomputer: 1,
		// },
		Name: "Supercomputer Alt",
		Input: map[Item]uint{
			Computer:                  3,
			ElectromagneticControlRod: 2,
			Battery:                   20,
			Wire:                      45,
		},
		Output: map[Item]uint{
			Supercomputer: 2,
		},
	},
	AutomatedWiring: {
		Name: "Automated Wiring Alt",
		Input: map[Item]uint{
			Stator:             2,
			Wire:               40,
			HighSpeedConnector: 1,
		},
		Output: map[Item]uint{
			AutomatedWiring: 4,
		},
	},
	HeavyModularFrame: {
		Name: "Heavy Modular Frame Alt",
		Input: map[Item]uint{
			ModularFrame:          8,
			EncasedIndustrialBeam: 10,
			SteelPipe:             36,
			Concrete:              22,
		},
		Output: map[Item]uint{
			HeavyModularFrame: 3,
		},
	},
	Battery: {
		Name: "Battery Alt",
		Input: map[Item]uint{
			Sulfur:              6,
			AlcladAluminumSheet: 7,
			Plastic:             8,
			Wire:                12,
		},
		Output: map[Item]uint{
			Battery: 4,
		},
	},
	ModularFrame: {
		Name: "Modular Frame Alt",
		Input: map[Item]uint{
			ReinforcedIronPlate: 2,
			SteelPipe:           10,
		},
		Output: map[Item]uint{
			ModularFrame: 3,
		},
	},
	Stator: {
		Name: "Stator Alt",
		Input: map[Item]uint{
			SteelPipe: 4,
			Quickwire: 15,
		},
		Output: map[Item]uint{
			Stator: 2,
		},
	},
	SteelPipe: {
		Name: "Steel Pipe",
		Input: map[Item]uint{
			SteelIngot: 3,
		},
		Output: map[Item]uint{
			SteelPipe: 2,
		},
	},
	SteelIngot: {
		Name: "Steel Ingot Alt",
		Input: map[Item]uint{
			IronIngot: 2,
			Coal:      2,
		},
		Output: map[Item]uint{
			SteelIngot: 3,
		},
	},
	Quickwire: {
		Name: "Quickwire Alt",
		Input: map[Item]uint{
			CateriumIngot: 1,
			CopperIngot:   5,
		},
		Output: map[Item]uint{
			Quickwire: 12,
		},
	},
	Computer: {
		Name: "Computer Alt",
		Input: map[Item]uint{
			CircuitBoard:      8,
			CrystalOscillator: 3,
		},
		Output: map[Item]uint{
			Computer: 3,
		},
	},
	ElectromagneticControlRod: {
		Name: "Electromagnetic Control Rod Alt",
		Input: map[Item]uint{
			Stator:             2,
			HighSpeedConnector: 1,
		},
		Output: map[Item]uint{
			ElectromagneticControlRod: 2,
		},
	},
	CateriumIngot: {
		Name: "Pure Caterium Ingot",
		Input: map[Item]uint{
			CateriumOre: 2,
			Water:       2,
		},
		Output: map[Item]uint{
			CateriumIngot: 1,
		},
	},
	Wire: {
		Name: "Wire Alt",
		Input: map[Item]uint{
			CopperIngot:   4,
			CateriumIngot: 1,
		},
		Output: map[Item]uint{
			Wire: 30,
		},
	},
	CopperIngot: {
		Name: "Pure Copper Ingot",
		Input: map[Item]uint{
			CopperOre: 6,
			Water:     4,
		},
		Output: map[Item]uint{
			CopperIngot: 15,
		},
	},
	Concrete: {
		Name: "Wet Concrete",
		Input: map[Item]uint{
			Limestone: 6,
			Water:     5,
		},
		Output: map[Item]uint{
			Concrete: 4,
		},
	},
	ReinforcedIronPlate: {
		Name: "Reinforced Plate Alt",
		Input: map[Item]uint{
			IronPlate: 10,
			Wire:      20,
		},
		Output: map[Item]uint{
			ReinforcedIronPlate: 3,
		},
	},
	CircuitBoard: {
		Name: "Circuit Board Alt",
		Input: map[Item]uint{
			CopperSheet: 11,
			Silica:      11,
		},
		Output: map[Item]uint{
			CircuitBoard: 5,
		},
	},
	Silica: {
		Name: "Silica Alt",
		Input: map[Item]uint{
			RawQuartz: 3,
			Limestone: 5,
		},
		Output: map[Item]uint{
			Silica: 7,
		},
	},
	HighSpeedConnector: {
		Name: "High-Speed Connector Alt",
		Input: map[Item]uint{
			Quickwire:    60,
			Silica:       25,
			CircuitBoard: 2,
		},
		Output: map[Item]uint{
			HighSpeedConnector: 2,
		},
	},
	CopperSheet: {
		Name: "Steamed Copper Sheet",
		Input: map[Item]uint{
			CopperIngot: 3,
			Water:       3,
		},
		Output: map[Item]uint{
			CopperSheet: 3,
		},
	},
	CrystalOscillator: {
		Name: "Crystal Oscillator Alt",
		Input: map[Item]uint{
			QuartzCrystal: 10,
			Rubber:        7,
			AILimiter:     1,
		},
		Output: map[Item]uint{
			CrystalOscillator: 1,
		},
	},
	AILimiter: {
		Name: "AI Limiter",
		Input: map[Item]uint{
			CopperSheet: 5,
			Quickwire:   20,
		},
		Output: map[Item]uint{
			AILimiter: 1,
		},
	},
	IronPlate: {
		Name: "Iron Plate Alt",
		Input: map[Item]uint{
			SteelIngot: 3,
			Plastic:    2,
		},
		Output: map[Item]uint{
			IronPlate: 18,
		},
	},
	EncasedIndustrialBeam: {
		Name: "Encased Industrial Pipe",
		Input: map[Item]uint{
			SteelPipe: 7,
			Concrete:  5,
		},
		Output: map[Item]uint{
			EncasedIndustrialBeam: 1,
		},
	},
	QuartzCrystal: {
		Name: "Quartz Crystal Alt",
		Input: map[Item]uint{
			RawQuartz: 9,
			Water:     5,
		},
		Output: map[Item]uint{
			QuartzCrystal: 7,
		},
	},
	AlcladAluminumSheet: {
		Name: "Alclad Aluminum Sheet",
		Input: map[Item]uint{
			AluminumIngot: 3,
			CopperIngot:   1,
		},
		Output: map[Item]uint{
			AlcladAluminumSheet: 3,
		},
	},
	MagneticFieldGenerator: {
		Name: "Magnetic Field Generator",
		Input: map[Item]uint{
			VersatileFramework:        5,
			ElectromagneticControlRod: 2,
			Battery:                   10,
		},
		Output: map[Item]uint{
			MagneticFieldGenerator: 2,
		},
	},
	VersatileFramework: {
		Name: "Versatile Framework Alt",
		Input: map[Item]uint{
			ModularFrame: 1,
			SteelBeam:    6,
			Rubber:       8,
		},
		Output: map[Item]uint{
			VersatileFramework: 2,
		},
	},
	SteelBeam: {
		Name: "Steel Beam",
		Input: map[Item]uint{
			SteelIngot: 4,
		},
		Output: map[Item]uint{
			SteelBeam: 1,
		},
	},
	Plastic: {
		Name: "Plastic Alt",
		Input: map[Item]uint{
			PolymerResin: 6,
			Water:        2,
		},
		Output: map[Item]uint{
			Plastic: 2,
		},
	},
	PolymerResin: {
		Name: "Polymer Resin Alt",
		Input: map[Item]uint{
			CrudeOil: 6,
		},
		Output: map[Item]uint{
			PolymerResin:    13,
			HeavyOilResidue: 2,
		},
	},
	Rubber: {
		Name: "Rubber Alt",
		Input: map[Item]uint{
			PolymerResin: 4,
			Water:        4,
		},
		Output: map[Item]uint{
			Rubber: 2,
		},
	},
	AluminumIngot: {
		
	}
}
