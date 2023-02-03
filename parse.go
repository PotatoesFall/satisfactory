package main

import (
	_ "embed"
	"encoding/json"

	"github.com/PotatoesFall/satisfactory/game"
)

//go:embed simple.json
var jsonFile []byte

func loadRecipes(config RecipeConfig) {
	var gameInfo game.Info
	if err := json.Unmarshal(jsonFile, &gameInfo); err != nil {
		panic(err)
	}
	allItems = gameInfo.Items

	recipesByName = make(map[string]*game.Recipe)
	recipesByItem = make(map[game.Item][]*game.Recipe)
outer:
	for _, recipe := range gameInfo.Recipes {
		for _, disallowed := range config.Disallowed {
			if disallowed == recipe.Name {
				continue outer
			}
		}

		recipesByName[recipe.Name] = recipe
		allRecipes = append(allRecipes, recipe)
		for item := range recipe.Products {
			recipesByItem[item] = append(recipesByItem[item], recipe)
		}
	}
}
