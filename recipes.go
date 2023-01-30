package main

import "encoding/json"

var allRecipes []Recipe

var recipesByItem map[Item][]Recipe

func loadRecipes() {
	if err := json.Unmarshal(recipesJSON, &allRecipes); err != nil {
		panic(err)
	}

	recipesByItem = make(map[Item][]Recipe)
	for _, recipe := range allRecipes {
		for item := range recipe.Output {
			recipesByItem[item] = append(recipesByItem[item], recipe)
		}
	}
}
