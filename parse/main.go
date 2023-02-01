package main

import (
	"encoding/json"
	"os"

	"github.com/PotatoesFall/satisfactory/game"
)

func main() {
	file, err := os.Open("game.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	classes := []struct {
		NativeClass string          `json:"NativeClass"`
		Classes     json.RawMessage `json:"Classes"`
	}{}

	if err := json.NewDecoder(file).Decode(&classes); err != nil {
		panic(err)
	}

	var resourcesJSON, itemsJSON, machinesJSON, recipesJSON, extractorJSON, particleAcceleratorJSON, biomassJSON, nuclearJSON []byte

	for _, class := range classes {
		switch class.NativeClass {
		case "Class'/Script/FactoryGame.FGResourceDescriptor'":
			resourcesJSON = class.Classes
		case "Class'/Script/FactoryGame.FGRecipe'":
			recipesJSON = class.Classes
		case "Class'/Script/FactoryGame.FGItemDescriptor'":
			itemsJSON = class.Classes
		case "Class'/Script/FactoryGame.FGBuildableManufacturer'":
			machinesJSON = class.Classes
		case "Class'/Script/FactoryGame.FGBuildableResourceExtractor'":
			extractorJSON = class.Classes
		case "Class'/Script/FactoryGame.FGBuildableManufacturerVariablePower'":
			particleAcceleratorJSON = class.Classes
		case "Class'/Script/FactoryGame.FGItemDescriptorBiomass'":
			biomassJSON = class.Classes
		case "Class'/Script/FactoryGame.FGItemDescriptorNuclearFuel'":
			nuclearJSON = class.Classes
			// case "Class'/Script/FactoryGame.FGBuildableGeneratorNuclear'":
			// 	generatorJSON = class.Classes
		}
	}

	machines, machineClasses := parseMachines(machinesJSON, particleAcceleratorJSON, extractorJSON)
	printMachines(machines)

	items, itemClasses := parseItems(itemsJSON, resourcesJSON, biomassJSON, nuclearJSON)
	printItems(items)

	recipes := parseRecipes(recipesJSON, machineClasses, itemClasses)
	recipes = append(recipes, &game.Recipe{
		Name: "Uranium Waste",
		Ingredients: map[game.Item]int{
			"Uranium Fuel Rod": 1,
			"Water":            1200,
		},
		Products: map[game.Item]int{
			"Uranium Waste": 50,
		},
		Duration: 300,
		Machine:  "Nuclear Power Plant",
		Power:    -2500,
	})
	recipes = append(recipes, &game.Recipe{
		Name: "Plutonium Waste",
		Ingredients: map[game.Item]int{
			"Plutonium Fuel Rod": 1,
			"Water":              2400,
		},
		Products: map[game.Item]int{
			"Plutonium Waste": 10,
		},
		Duration: 600,
		Machine:  "Nuclear Power Plant",
		Power:    -2500,
	})
	printRecipes(recipes)

	gameFile, err := os.Create("simple.json")
	if err != nil {
		panic(err)
	}
	defer gameFile.Close()

	enc := json.NewEncoder(gameFile)
	enc.SetIndent("", " ")
	err = enc.Encode(game.Info{
		Items:   items,
		Recipes: recipes,
	})
	if err != nil {
		panic(err)
	}
}
