package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/PotatoesFall/satisfactory/game"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"golang.org/x/exp/slices"
)

/*
TODO
- allow for multiple inputs
- account for byproducts
	- combine with other recipes
		- in ratios if possible, that would be dank
		- idea, build macro recipes during the getRecipeWeights phase?
			- a different kind of recipe
			- contains a list of sub recipes
			- no byproducts
				- consumed in tree
					- partially consumed? (then sink the rest)
				- burned for energy
	- generate energy if possible
	- process and sink if liquid ?
- add energy costs for resource extraction
- add building cost
	- also for resoure extraction
*/

type Config struct {
	Weights Weights      `toml:"Weights"`
	Recipes RecipeConfig `toml:"Recipes"`
}

type RecipeConfig struct {
	Disallowed []string `toml:"Disallowed"`
}

type Weights struct {
	Resources map[game.Item]float64 `toml:"Resources"`

	Power float64 `toml:"Power"`

	// Machines map[string]float64
}

var (
	allRecipes    []*game.Recipe
	allItems      []game.Item
	recipesByItem map[game.Item][]*game.Recipe
	recipesByName map[string]*game.Recipe
)

func main() {
	config := readConfig()
	loadRecipes(config.Recipes)

	item, amount := getItem()
	_, _ = item, amount

	recipeCounts := optimize(allRecipes, config.Weights, map[game.Item]float64{item: amount})
	// recipeCounts := optimize([]*game.Recipe{
	// 	recipesByName["Alternate: Pure Iron Ingot"],
	// 	recipesByName["Iron Ingot"],
	// }, config.Weights, map[game.Item]float64{"Iron Ingot": 60})

	total := 0.0
	for recipe, count := range recipeCounts {
		power := recipe.Power * count
		total += power
		fmt.Println(count, recipe.Name, power, "MW")
	}
	fmt.Println("Total", total, "MW")

	// recipeTrees := getAllItemWeights(weights)
	// tree := recipeTrees[item]

	// // fmt.Println(tree.Print(amount))

	// recipeOrder := tree.RecipeOrder()
	// fmt.Println(markdownTable(tree.RecipeCounts(amount), recipeOrder))

	// fmt.Printf("Total Power: %.2f MW", tree.Power(amount))
}

//go:embed config.example.toml
var defaultConfig []byte

func readConfig() Config {
	var config Config

	configFile, err := os.ReadFile("config.toml")
	if err != nil {
		fmt.Println("No config.toml found, using default config.")
		configFile = defaultConfig
	}

	if err := toml.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}

	return config
}

func markdownTable(recipeCounts map[string]float64, recipeOrder []string) string {
	rows := make([][]string, 0, len(recipeCounts))
	for _, recipeName := range recipeOrder {
		count := recipeCounts[recipeName]
		recipe, ok := recipesByName[recipeName]
		if !ok {
			rows = append(rows, []string{
				"", "", "", recipeName + " (Resource)", fmt.Sprintf("%8.2f %25s", count, recipeName), "",
			})
			continue
		}

		var ingredients, products strings.Builder
		for ingredient, inCount := range recipe.Ingredients {
			ingredients.WriteString(fmt.Sprintf("%8.2f %25s", float64(inCount)*count*60/recipe.Duration, ingredient))
		}
		for product, prodCount := range recipe.Products {
			products.WriteString(fmt.Sprintf("%8.2f %25s", float64(prodCount)*count*60/recipe.Duration, product))
		}

		rows = append(rows, []string{
			fmt.Sprintf("%7.2f", count),
			recipe.Machine,
			fmt.Sprintf("%7.2f MW", recipe.Power*count),
			recipeName,
			products.String(),
			ingredients.String(),
		})
	}

	builder := markdown.NewTableFormatterBuilder().WithPrettyPrint()
	f := builder.Build("", "Machine", "Power", "Recipe", "Products", "Ingredients")
	table, err := f.Format(rows)
	if err != nil {
		panic(err)
	}

	return table
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
