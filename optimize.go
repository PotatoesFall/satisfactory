package main

import (
	"fmt"
	"math"

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

/*
NOTES AFTER MUCH DOING
- phases
	- startup phase
		- have a global factor that starts small and slowly increases to avoid jumps
		- a high weight for products
		- a low (no?) weight for ingredients (?)
		- a low (no?) weight for power and resources
	- balancing phase
		- equal weight for ingredients and products (?)
		- higher weight for resources and power
		- steadily increase global factor ? or just ingredients and products?
	- correction phase
		- steadliy increase relative weight of products and ingredients to get closer to perfect math
another idea: concurrency for performance?
*/

type OptimizationWeights struct {
	Global float64

	Resources             float64 // also use for power?
	Ingredients, Products float64
}

const (
	maxRecipeAmount = 1e3

	startupRounds    = 1_000
	balancingRounds  = 100_000
	correctingRounds = 10_000
	eliminationStart = startupRounds

	// computed
	totalRounds     = startupRounds + balancingRounds + correctingRounds
	correctingStart = startupRounds + balancingRounds + 1

	// global factors
	globalFactor      = 1e-14
	resourceFactor    = 1.0
	correctnessFactor = 1e10

	// startup
	startupFactor            = 1e0
	startupResourceFactor    = 1.0
	startupIngredientsFactor = 1e-3

	// balancing
	balancingFactor            = 1.0
	balancingScaling           = 1e1
	balancingResourceFactor    = 1e8
	balancingIngredientsFactor = 1.0

	// correcting
	correctingFactor            = 1e5
	correctingResourceFactor    = 1.0
	correctingIngredientsFactor = 1.0

	// recipe Amounts
	initialAmount     = 0.0
	eliminationAmount = 1e-12

	// logging
	log       = true
	logRecipe = "Packaged Fuel"
	// logEntries  = 1000
	stopRound   = 10_000
	logEntries  = 10_000
	logInterval = totalRounds / logEntries
)

func getOptWeights(round int) OptimizationWeights {
	w := OptimizationWeights{
		Global:      globalFactor,
		Resources:   resourceFactor,
		Ingredients: correctnessFactor,
		Products:    correctnessFactor,
	}

	switch {
	case round < startupRounds:
		// startupProgress := float64(round) / startupRounds
		w.Global *= startupFactor
		w.Resources *= startupResourceFactor
		w.Ingredients *= startupIngredientsFactor

	case round-startupRounds < balancingRounds:
		// progress := float64(round-startupRounds) / balancingRounds
		progress := math.Exp(balancingScaling * float64(round-startupRounds) / balancingRounds)
		w.Global *= progress * balancingFactor
		w.Resources *= balancingResourceFactor
		w.Ingredients *= balancingIngredientsFactor

	default: // correcting round
		progress := float64(round-startupRounds) / correctingRounds
		w.Global *= progress * correctingFactor
		w.Resources *= correctingResourceFactor
		w.Ingredients *= correctingIngredientsFactor
	}

	return w
}

func optimize(allRecipes []*game.Recipe, w Weights, requirements map[game.Item]float64) map[*game.Recipe]float64 {
	recipeCounts := make(map[*game.Recipe]float64)
	for _, recipe := range allRecipes {
		recipeCounts[recipe] = 0
	}
	derivatives := make(map[*game.Recipe]float64, len(recipeCounts))

	for i := 0; i < totalRounds; i++ {
		if i > stopRound {
			break
		}

		fluxes := sumRecipes(recipeCounts)
		for req, flux := range requirements {
			fluxes[req] -= flux
		}

		// logging
		if log && i%logInterval == 0 {
			fmt.Println("ROUND", i)
			for _, recipe := range allRecipes {
				count := recipeCounts[recipe]
				if count > eliminationAmount {
					fmt.Printf("\t%.5f %s\n", count, recipe.Name)
				}
			}
			for item, flux := range fluxes {
				fmt.Println("\t\tFLUX", item, flux)
			}
		}

		// get weights based on round
		optWeights := getOptWeights(i)

		// get all derivatives
		for recipe := range recipeCounts {
			derivatives[recipe] = getDerivative(w, recipe, fluxes, optWeights)

			// logging
			if log && recipe.Name == logRecipe {
				fmt.Println(" round", i, "d/dx:", derivatives[recipe], recipe.Name, "count:", recipeCounts[recipe])
			}
		}

		// apply derivative
		for recipe, derivative := range derivatives {
			count := recipeCounts[recipe]

			newCount := math.Min(math.Max(0, count-derivative*optWeights.Global), maxRecipeAmount)

			recipeCounts[recipe] = newCount

			if i > eliminationStart && newCount <= eliminationAmount {
				delete(recipeCounts, recipe)
			}
		}
	}

	return recipeCounts
}

func getDerivative(w Weights, recipe *game.Recipe, fluxes map[game.Item]float64, optWeight OptimizationWeights) float64 {
	powerDerivative := recipe.Power * w.Power * optWeight.Resources

	ingredientDerivative := 0.0
	for item, count := range recipe.Ingredients {
		if weight, ok := w.Resources[item]; ok {
			ingredientDerivative += optWeight.Resources * weight * float64(count) / recipe.DurationMinutes()
			continue
		}
		if flux := fluxes[item]; flux < 0 {
			ingredientDerivative += optWeight.Ingredients * (-flux) * float64(count) / recipe.DurationMinutes()
		}
	}

	productDerivative := 0.0
	for item, count := range recipe.Products {
		if weight, ok := w.Resources[item]; ok {
			if fluxes[item] <= 0 {
				productDerivative -= optWeight.Resources * weight * float64(count) / recipe.DurationMinutes()
			}
			continue
		}

		if flux := fluxes[item]; flux < 0 {
			productDerivative -= optWeight.Products * (-flux) * float64(count) / recipe.DurationMinutes()
		}
	}

	if log && recipe.Name == logRecipe {
		fmt.Print("DERIVATIVE ", "power: ", powerDerivative, " ingredient: ", ingredientDerivative, " product: ", productDerivative)
	}

	return powerDerivative + ingredientDerivative + productDerivative
}

func sumRecipes(recipeCounts map[*game.Recipe]float64) (fluxes map[game.Item]float64) {
	fluxes = make(map[game.Item]float64, len(allItems))

	for recipe, recipeAmount := range recipeCounts {
		if recipeAmount <= eliminationAmount {
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
