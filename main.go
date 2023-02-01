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
			"Iron Ore":     1, // https://satisfactory.fandom.com/wiki/Resource_node
			"Copper Ore":   70380.0 / 28860,
			"Limestone":    70380.0 / 52860,
			"Water":        0,
			"Sulfur":       70380.0 / 6840,
			"Coal":         70380.0 / 30900,
			"Raw Quartz":   70380.0 / 10500,
			"Caterium Ore": 70380.0 / 11040,
			"Crude Oil":    70380.0 / 9900,
			"Nitrogen Gas": 70380.0 / 10000, // guess, not on the wiki
			"Uranium":      70380.0 / 2100,
			"Bauxite":      70380.0 / 9780,

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

	fmt.Println("|       | Building     | Power        | Recipe                                |")
	fmt.Println("| ----- | -------------|--------------|---------------------------------------|")
	for recipeName, count := range tree.RecipeCounts(amount) {
		recipe := recipesByName[recipeName]
		fmt.Printf("|%6.2fx|%14s|%11s MW|%39s|\n", count, recipe.Machine, fmtAmount(recipe.Power*count), recipeName)
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
