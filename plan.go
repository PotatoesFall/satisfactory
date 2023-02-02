package main

import (
	"math"

	"github.com/PotatoesFall/satisfactory/game"
)

type Plan interface {
	Ingredients() map[game.Item]float64
	Products() map[game.Item]float64
	Power() float64
	RecipeCounts() map[string]float64
	Cost(Weights) float64
}

type simpleRecipePlan struct {
	Recipe *game.Recipe
}

func (srp simpleRecipePlan) Ingredients() map[game.Item]float64 {
	fluxes := make(map[game.Item]float64)
	for ingredient, count := range srp.Recipe.Ingredients {
		fluxes[ingredient] = float64(count) * srp.Recipe.RatePerMinute()
	}
	return fluxes
}

func (srp simpleRecipePlan) Products() map[game.Item]float64 {
	fluxes := make(map[game.Item]float64)
	for product, count := range srp.Recipe.Products {
		fluxes[product] = float64(count) * srp.Recipe.RatePerMinute()
	}
	return fluxes
}

func (srp simpleRecipePlan) Power() float64 {
	return srp.Recipe.Power
}

func (srp simpleRecipePlan) RecipeCounts() map[string]float64 {
	return map[string]float64{
		srp.Recipe.Name: 1,
	}
}

func (srp simpleRecipePlan) Cost(w Weights) float64 {
	return srp.Recipe.Power * w.Power
}

type resourcePlan struct {
	Item game.Item
}

func (rp resourcePlan) Ingredients() map[game.Item]float64 {
	return nil
}

func (rp resourcePlan) Products() map[game.Item]float64 {
	return map[game.Item]float64{
		rp.Item: 1,
	}
}

func (rp resourcePlan) Power() float64 {
	return 0 // TODO
}

func (rp resourcePlan) RecipeCounts() map[string]float64 {
	return map[string]float64{
		string(rp.Item) + " (Resource)": 1,
	}
}

func (rp resourcePlan) Cost(w Weights) float64 {
	return w.Resources[rp.Item]
}

type scaledPlan struct {
	Plan       Plan
	Multiplier float64
}

func (sp scaledPlan) Ingredients() map[game.Item]float64 {
	ingredients := make(map[game.Item]float64)
	for i, f := range sp.Plan.Ingredients() {
		ingredients[i] = f * sp.Multiplier
	}
	return ingredients
}

func (sp scaledPlan) Products() map[game.Item]float64 {
	products := make(map[game.Item]float64)
	for i, f := range sp.Plan.Products() {
		products[i] = f * sp.Multiplier
	}
	return products
}

func (sp scaledPlan) Power() float64 {
	return sp.Plan.Power() * sp.Multiplier
}

func (sp scaledPlan) RecipeCounts() map[string]float64 {
	counts := make(map[string]float64)
	for i, f := range sp.Plan.RecipeCounts() {
		counts[i] = f * sp.Multiplier
	}
	return counts
}

func (sp scaledPlan) Cost(w Weights) float64 {
	return sp.Plan.Cost(w) * sp.Multiplier
}

type multiPlan []Plan

func (mp multiPlan) Ingredients() map[game.Item]float64 {
	ingredients := make(map[game.Item]float64)
	for _, plan := range mp {
		for i, f := range plan.Ingredients() {
			ingredients[i] = f
		}
	}
	for _, plan := range mp {
		for i, f := range plan.Products() {
			ingredients[i] = math.Max(0, ingredients[i]-f)
		}
	}
	return ingredients
}

func (mp multiPlan) Products() map[game.Item]float64 {
	products := make(map[game.Item]float64)
	for _, plan := range mp {
		for i, f := range plan.Products() {
			products[i] = f
		}
	}
	for _, plan := range mp {
		for i, f := range plan.Ingredients() {
			products[i] = math.Max(0, products[i]-f)
		}
	}
	return products
}

func (mp multiPlan) Power() float64 {
	power := 0.0
	for _, plan := range mp {
		power += plan.Power()
	}
	return power
}

func (mp multiPlan) RecipeCounts() map[string]float64 {
	counts := make(map[string]float64)
	for _, plan := range mp {
		for recipe, count := range plan.RecipeCounts() {
			counts[recipe] += count
		}
	}
	return counts
}

func (mp multiPlan) Cost(w Weights) float64 {
	cost := 0.0
	for _, plan := range mp {
		cost += plan.Cost(w)
	}
	return cost
}

// assert interface
var (
	_ Plan = simpleRecipePlan{}
	_ Plan = resourcePlan{}
	_ Plan = scaledPlan{}
	_ Plan = multiPlan{}
	// TODO: cachedPlan?
)
