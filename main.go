package main

import (
	"fmt"

	_ "embed"
)

/*
TODO
- make the system be able to check different recipe variants
	- add weights for
		- resource rarity (per resource)
		- energy cost (don't forget miners/extractors)
		- number of machines? (per machine) nice-to-haven
- have the system account for by-products
	- match by-product recipes and balance them perfectly (how?)
	- nice-to-have: account for energy that can be produced with byproducts and deduct from total energy cost
- allow choosing/locking certain recipes while letting others be optimized?
*/

//go:embed recipes.json
var recipesJSON []byte

type Weights struct {
	Ore map[Item]float64

	Energy float64

	// Machines float64
}

func main() {
	loadRecipes()

	// TODO: make parameters
	weights := Weights{
		Ore: map[Item]float64{
			IronOre:     1,
			CopperOre:   1,
			Limestone:   1,
			Water:       1,
			Sulfur:      1,
			Coal:        1,
			RawQuartz:   1,
			CateriumOre: 1,
			CrudeOil:    1,
		},
		Energy: 0,
	}

	recipeTrees := getAllItemWeights(weights)

	adsTree := recipeTrees[AssemblyDirectorSystem]

	fmt.Println(adsTree.Print(4))

	recipeCounts := adsTree.RecipeCounts(4)
	for recipeName, count := range recipeCounts {
		recipe := recipesByName[recipeName]
		fmt.Printf("%6.2f %s (%s) %s MW\n", count, recipe.Machine, recipeName, fmtAmount(recipe.Power*count))
	}
}
