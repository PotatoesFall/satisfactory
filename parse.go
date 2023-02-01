package main

import (
	_ "embed"
	"encoding/json"

	"github.com/PotatoesFall/satisfactory/game"
)

//go:embed simple.json
var jsonFile []byte

func loadRecipes() {
	var gameInfo game.Info
	if err := json.Unmarshal(jsonFile, &gameInfo); err != nil {
		panic(err)
	}
	allRecipes = gameInfo.Recipes
	allItems = gameInfo.Items

	recipesByName = make(map[string]*game.Recipe)
	recipesByItem = make(map[game.Item][]*game.Recipe)
	for _, recipe := range gameInfo.Recipes {
		recipesByName[recipe.Name] = recipe
		for item := range recipe.Products {
			recipesByItem[item] = append(recipesByItem[item], recipe)
		}
	}
}
