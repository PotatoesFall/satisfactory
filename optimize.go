package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/PotatoesFall/satisfactory/game"
)

/*
STRATEGY
to get best plan for all outcomes
	- make a map of recipe to amount
	- add all recipes for making any of the outcomes
		- get recipes for one outcome
			- if the outcome is a resource, done.
			- get all recipes that produce this item, excluding recipes already used in the tree to avoid infinite trees
			- RECURSE: for each recipe's ingredients, find all viable recipes
		- add them to the set
	- balance all recipes against each other (regression? elimination?)
		- consider a cacheLayer around each recipe?
		- start with all recipes in there
		- LOOP x times:
			- scale up/down equally as needed
			- calculate initial cost
			- for each recipe, calculate a derivative
				- make a copy of the current recipes
				- scale up the plan in question by an infinitesmal amount
				- scale all plans up/down equally as needed
				- calculate new cost and determine derivative
			- adjust all plans according to their derivative
				- the amount should slowly decrease over time
			- if one plan gets below a certain threshold, eliminate that plan completely, adjusting the rest as needed.
		- at the end, stop optimizing ratios and just scale down each recipe as much as possible, going right to left, to get the math exactly correct.
			- or maybe don't, just keep it simple
		- round and use final ratios!


okay Im liking the "all recipes" approach it just has some problems:
	- producing not enough input or using resources that aren't there needs to be added to the cost, and should scale non-linearly
	- if we do that we don't need to scale up and down each round
	- we can optimize the derivative function by calculating the inputs, outputs and cost once, and then durung the derivative function we don't need to recalculate them
		- output differentials - how far are we from the desired outcome
			- punish the more we go below, do nothing if we go above
			- edit: this should be e^x or whatever so we can calculate the ACTUAL derivative
		- input differentials - are we using more than we produce
			- same as above
			- for resources just use the weight duh, but don't forget them!
		- don't forget about power
*/

const (
	log = false
	// logRecipe = ""
	logRecipe = "iAlternate: Cast Screw"

	filter       = false
	initialCount = minValue
	historyLimit = 7

	startFactor = 10000000
	factor      = 10

	rounds        = 100_000
	startupRounds = 1000
	// logInterval = rounds / 1000
	logInterval = 1
	// eliminationStart = rounds - 100 // off until last couple rounds
	eliminationStart = rounds / 10

	minValue         = 0.00000000000000000000000000000001
	derivativeFactor = 1

	productFactor    = 0.001
	ingredientFactor = productFactor / 1000
)

func filterRecipes(recipes []*game.Recipe, requirements map[game.Item]float64) map[*game.Recipe]float64 {
	counts := map[*game.Recipe]float64{}

	if filter {
		for req, flux := range requirements {
			_filterRecipes(counts, recipes, req, flux, nil)
		}
	}

	for _, recipe := range recipes {
		if counts[recipe] == 0 {
			counts[recipe] = initialCount
		}

		// TODO fix this issue
		if strings.Contains(recipe.Name, "ackage") {
			delete(counts, recipe)
		}
	}

	return counts
}

func _filterRecipes(counts map[*game.Recipe]float64, recipes []*game.Recipe, req game.Item, flux float64, history []*game.Recipe) {
	if len(history) > historyLimit {
		return
	}
	var relevantRecipes []*game.Recipe

outer:
	for _, recipe := range recipes {
		for _, hist := range history {
			if hist == recipe {
				continue outer
			}
		}

		for product := range recipe.Products {
			if product == req {
				relevantRecipes = append(relevantRecipes, recipe)
				continue
			}
		}
	}

	recipeCount := len(relevantRecipes)
	for _, recipe := range relevantRecipes {
		production := float64(recipe.Products[req]) / recipe.DurationMinutes()
		recipeRate := flux / production
		counts[recipe] += recipeRate / float64(recipeCount)
		newHistory := append(history, recipe)
		for ingredient, inRatio := range recipe.Ingredients {
			ingredientFlux := recipeRate * float64(inRatio)
			_filterRecipes(counts, recipes, ingredient, ingredientFlux, newHistory)
		}
	}
}

func optimize(allRecipes []*game.Recipe, w Weights, requirements map[game.Item]float64) map[*game.Recipe]float64 {
	recipeCounts := filterRecipes(allRecipes, requirements)
	derivatives := make(map[*game.Recipe]float64, len(recipeCounts))

	for i := 0; i < rounds; i++ {
		roundFactor := startFactor + float64(i)*factor
		if i >= rounds*9/10 {
			// second half --> lets go
			roundFactor *= math.Pow(1.005, float64(i-rounds*9/10)/10)
			// fmt.Println(roundFactor)
		}

		// fmt.Println(roundFactor)
		fluxes := sumRecipes(recipeCounts)
		for req, flux := range requirements {
			fluxes[req] -= flux
		}

		// logging
		if log && i%logInterval == 0 {
			fmt.Println("ROUND", i)
			for _, recipe := range allRecipes {
				count := recipeCounts[recipe]
				if count > minValue {
					fmt.Printf("\t%.5f %s\n", count, recipe.Name)
				}
			}
			for item, flux := range fluxes {
				fmt.Println("\t\tFLUX", item, flux)
			}
		}

		// get all derivatives
		for recipe := range recipeCounts {
			// logging
			if recipe.Name == logRecipe {
				fmt.Print("DERIVATIVE: ", recipe.Name, " ", recipeCounts[recipe], " ")
			}

			derivatives[recipe] = getDerivative(w, recipe, fluxes, roundFactor)
		}

		// apply derivative
		for recipe, derivative := range derivatives {
			count := recipeCounts[recipe]

			// delta := -math.Max(
			// 	math.Min(
			// 		derivative*derivativeFactor/factor,
			// 		maxDecrement/factor,
			// 	),
			// 	-maxIncrement/factor)
			delta := -derivative * derivativeFactor / roundFactor
			if i < startupRounds*2 {
				if i < startupRounds {
					delta *= 0.1
				}
				delta *= 0.1
			}

			recipeCounts[recipe] = math.Max(
				count+delta,
				minValue)

			if i > eliminationStart && recipeCounts[recipe] <= minValue {
				delete(recipeCounts, recipe)
			}
		}
	}

	return recipeCounts
}

func getDerivative(w Weights, recipe *game.Recipe, fluxes map[game.Item]float64, factor float64) float64 {
	powerDerivative := recipe.Power * w.Power

	ingredientDerivative := 0.0
	for item, count := range recipe.Ingredients {
		if weight, ok := w.Resources[item]; ok {
			ingredientDerivative += weight * float64(count) / recipe.DurationMinutes()
			continue
		}
		if flux := fluxes[item]; flux < 0 {
			ingredientDerivative += ingredientFactor * (-flux) * float64(count) / recipe.DurationMinutes() * factor
		}
	}

	productDerivative := 0.0
	for item, count := range recipe.Products {
		if weight, ok := w.Resources[item]; ok {
			if fluxes[item] <= 0 {
				productDerivative -= weight * float64(count) / recipe.DurationMinutes()
			}
			continue
		}

		if flux := fluxes[item]; flux < 0 {
			productDerivative -= productFactor * (-flux) * float64(count) / recipe.DurationMinutes() * factor
		}
	}

	if recipe.Name == logRecipe {
		fmt.Println("power", powerDerivative, "ingredient", ingredientDerivative, "product", productDerivative)
	}

	return powerDerivative + ingredientDerivative + productDerivative
}

func sumRecipes(recipeCounts map[*game.Recipe]float64) (fluxes map[game.Item]float64) {
	fluxes = make(map[game.Item]float64, len(allItems))

	for recipe, recipeAmount := range recipeCounts {
		if recipeAmount <= minValue {
			continue
		}
		for input, amount := range recipe.Ingredients {
			fluxes[input] -= float64(amount) * recipeAmount / recipe.DurationMinutes()
		}
		for input, amount := range recipe.Products {
			fluxes[input] += float64(amount) * recipeAmount / recipe.DurationMinutes()
		}
	}

	return fluxes
}
