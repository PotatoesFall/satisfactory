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
		if weight, ok := w.Resources[item]; ok {
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

			weight := w.Power * recipe.Power

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

	rt.traverse(flux, 0, func(tree *RecipeTree, flux float64, _ int) {
		if tree.Recipe == nil {
			recipes[string(tree.Item)] += flux
			return
		}

		recipeCount := flux / float64(tree.Recipe.Products[tree.Item]) * tree.Recipe.Duration / 60
		recipes[tree.Recipe.Name] += recipeCount
	})

	return recipes
}

func (rt *RecipeTree) RecipeOrder() []string {
	seenItems := make(map[game.Item]bool)
	var order []string
	seenRecipes := make(map[string]bool)

	recipeCounts := rt.RecipeCounts(1)

	for len(order) != len(recipeCounts) {
		newSeenItems := make(map[game.Item]bool)
	outer:
		for recipeName := range recipeCounts {
			if seenRecipes[recipeName] {
				continue
			}
			recipe, ok := recipesByName[recipeName]
			if !ok { // raw resource
				seenRecipes[recipeName] = true
				newSeenItems[game.Item(recipeName)] = true
				order = append(order, recipeName)
				continue
			}

			for input := range recipe.Ingredients {
				if !seenItems[input] { // only add if all ingredients have been added
					continue outer
				}
			}

			order = append(order, recipeName)
			for product := range recipe.Products {
				newSeenItems[product] = true
			}
			seenRecipes[recipeName] = true
		}

		for item := range newSeenItems {
			seenItems[item] = true
		}
	}

	// reverse order
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}

func (rt *RecipeTree) Power(flux float64) float64 {
	power := 0.0

	rt.traverse(flux, 0, func(tree *RecipeTree, flux float64, _ int) {
		if tree.Recipe == nil {
			return
		}

		power += flux / float64(tree.Recipe.Products[tree.Item]) * tree.Recipe.Power * tree.Recipe.Duration / 60
	})

	return power
}

func (rt *RecipeTree) traverse(flux float64, depth int, f func(tree *RecipeTree, flux float64, _ int)) {
	f(rt, flux, depth)

	for _, input := range rt.Inputs {
		ratio := flux * float64(rt.Recipe.Ingredients[input.Item]) / float64(rt.Recipe.Products[rt.Item])
		input.traverse(ratio, depth, f)
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
