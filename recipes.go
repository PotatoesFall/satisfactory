package main

import "encoding/json"

var (
	allRecipes    []*Recipe
	recipesByItem map[Item][]*Recipe
	recipesByName map[string]*Recipe
)

func loadRecipes() {
	if err := json.Unmarshal(recipesJSON, &allRecipes); err != nil {
		panic(err)
	}

	recipesByName = make(map[string]*Recipe)
	recipesByItem = make(map[Item][]*Recipe)
	for _, recipe := range allRecipes {
		recipesByName[recipe.Name] = recipe
		for item := range recipe.Output {
			recipesByItem[item] = append(recipesByItem[item], recipe)
		}
	}

}
