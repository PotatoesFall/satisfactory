package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PotatoesFall/satisfactory/game"
)

func getAllItemWeights(w Weights) map[game.Item]*RecipeTree {
	trees := make(map[game.Item]*RecipeTree)
	for _, item := range allItems {
		if weight, ok := w.Base[item]; ok {
			trees[item] = &RecipeTree{
				Item:   item,
				Weight: weight,
			}
		}
	}

	recipesChecked := make(map[string]bool)

	for len(recipesChecked) < len(allRecipes) {
	outer:
		for _, recipe := range allRecipes {
			if recipesChecked[recipe.Name] {
				continue
			}

			var weight float64
			inputs := make([]*RecipeTree, 0, len(recipe.Ingredients))
			for item, amount := range recipe.Ingredients {
				// only check recipes with all known ingredients
				tree, ok := trees[item]
				if !ok {
					continue outer
				}

				weight += tree.Weight * float64(amount)
				inputs = append(inputs, tree)
			}

			for item, amount := range recipe.Products {
				outWeight := weight / float64(amount)
				if tree, ok := trees[item]; ok && outWeight > tree.Weight {
					continue
				}

				trees[item] = &RecipeTree{
					Item:   item,
					Weight: outWeight,
					Recipe: recipe,
					Inputs: inputs,
				}
			}

			recipesChecked[recipe.Name] = true
		}
	}

	return trees
}

type RecipeTree struct {
	Item   game.Item // flux is always one
	Weight float64

	Recipe *game.Recipe
	Inputs []*RecipeTree
}

func (rt *RecipeTree) Print(flux float64) string {
	var buf strings.Builder
	printBuildTree(&buf, rt, flux, 0)
	return buf.String()
}

func (rt *RecipeTree) RecipeCounts(flux float64) map[string]float64 {
	recipes := make(map[string]float64)

	getRecipeCounts(recipes, rt, flux)

	return recipes
}

func getRecipeCounts(recipes map[string]float64, tree *RecipeTree, flux float64) {
	if tree.Recipe == nil {
		return
	}

	recipes[tree.Recipe.Name] += flux / float64(tree.Recipe.Products[tree.Item])

	for _, input := range tree.Inputs {
		ratio := flux * float64(tree.Recipe.Ingredients[input.Item]) / float64(tree.Recipe.Products[tree.Item])
		getRecipeCounts(recipes, input, ratio)
	}
}

func (rt *RecipeTree) Resources() map[game.Item]float64 {
	panic("not implemented")
}

func printBuildTree(buf *strings.Builder, tree *RecipeTree, flux float64, indentation int) {
	prefix := '-'
	if _, ok := recipesByItem[tree.Item]; !ok {
		prefix = '■'
	}

	recipeName := "Resource"
	if tree.Recipe != nil {
		recipeName = tree.Recipe.Name
	}

	buf.WriteString(fmt.Sprintf("%s%c %s %s (%s)\n", strings.Repeat("\t", indentation), prefix, fmtAmount(flux), tree.Item, recipeName))
	for _, input := range tree.Inputs {
		ratio := flux * float64(tree.Recipe.Ingredients[input.Item]) / float64(tree.Recipe.Products[tree.Item])
		printBuildTree(buf, input, ratio, indentation+1)
	}
}

func fmtAmount(amount float64) string {
	str := strconv.FormatFloat(amount, 'f', 5, 64)
	var i int
	for i = len(str) - 1; i >= 0; i-- {
		if str[i] == '.' {
			break
		}
		if str[i] != '0' {
			i++
			break
		}
	}
	return str[0:i]
}
