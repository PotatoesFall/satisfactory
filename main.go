package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PotatoesFall/satisfactory/game"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"golang.org/x/exp/slices"
)

/*
TODO
- add energy costs
	- also for resource extraction
- add building cost
	- also for resoure extraction
- account for byproducts
	- generate energy if possible
	- process and sink if liquid ?
*/

type Weights struct {
	Base map[game.Item]float64

	// Energy float64

	// Machines map[string]float64
}

var (
	allRecipes    []*game.Recipe
	allItems      []game.Item
	recipesByItem map[game.Item][]*game.Recipe
	recipesByName map[string]*game.Recipe
)

func main() {
	loadRecipes()
	item, amount := getItem()

	weights := Weights{
		Base: map[game.Item]float64{
			"Iron Ore":     70380, // https://satisfactory.fandom.com/wiki/Resource_node
			"Copper Ore":   28860 / 70380.0,
			"Limestone":    52860 / 70380.0,
			"Water":        0,
			"Sulfur":       6840 / 70380.0,
			"Coal":         30900 / 70380.0,
			"Raw Quartz":   10500 / 70380.0,
			"Caterium Ore": 11040 / 70380.0,
			"Crude Oil":    9900 / 70380.0,
			"Nitrogen Gas": 10000 / 70380.0, // guess, not on the wiki
			"Uranium":      2100 / 70380.0,
			"Bauxite":      9780 / 70380.0,

			"Yellow Power Slug":      100_000,
			"Purple Power Slug":      100_000,
			"Blue Power Slug":        100_000,
			"Leaves":                 100_000,
			"Wood":                   100_000,
			"Plasma Spitter Remains": 100_000,
			"Hog Remains":            100_000,
			"Stinger Remains":        100_000,
			"Mycelia":                100_000,
			"FICSMAS Gift":           100_000,
			"Flower Petals":          100_000,
			"Hatcher Remains":        100_000,
		},

		// Energy: 1,
	}

	recipeTrees := getAllItemWeights(weights)

	tree := recipeTrees[item]

	// fmt.Println(tree.Print(amount))

	fmt.Println("\nRECIPES")
	fmt.Println("|       | Building      | Power       | Recipe                            |")
	fmt.Println("| ----- | --------------|-------------|-----------------------------------|")
	for recipeName, count := range tree.RecipeCounts(amount) {
		recipe := recipesByName[recipeName]
		fmt.Printf("|%6.2fx|%15s|%10s MW|%35s|\n", count, recipe.Machine, fmtAmount(recipe.Power*count), recipeName)
	}

	// fmt.Println("\nRESOURCES")
	// for resource, count := range tree.Resources() {
	// 	fmt.Printf("%.2f %s", count, resource)
	// }
}

func getItem() (game.Item, float64) {
	if len(os.Args) < 3 {
		return pickItem()
	}

	f, err := strconv.ParseFloat(os.Args[len(os.Args)-1], 64)
	if err != nil {
		panic(err)
	}

	item := game.Item(strings.Join(os.Args[1:len(os.Args)-1], " "))
	if _, ok := recipesByItem[item]; !ok {
		panic("item not found: " + item)
	}

	return item, f
}

func pickItem() (game.Item, float64) {
	slices.Sort(allItems)
	itemPicker := selection.New("Produce What?", allItems)
	itemPicker.PageSize = 10

	item, err := itemPicker.RunPrompt()
	if err != nil {
		panic(err)
	}

	amountPicker := textinput.New("How much?")
	amountPicker.Placeholder = "2.5"
	amountPicker.Validate = func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		return err
	}

	s, err := amountPicker.RunPrompt()
	if err != nil {
		panic(err)
	}
	amount, _ := strconv.ParseFloat(s, 64)

	return item, amount
}
