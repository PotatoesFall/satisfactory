package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/PotatoesFall/satisfactory/game"
)

func parseRecipes(data []byte, machines map[string]RawMachine, items map[string]game.Item) []*game.Recipe {
	var rawRecipes []RawRecipe
	if err := json.Unmarshal(data, &rawRecipes); err != nil {
		panic(err)
	}

	recipeNames := map[string]int{}

	var recipes []*game.Recipe
	for _, raw := range rawRecipes {
		machine, ok := machines[string(raw.Machine)]
		if !ok {
			continue
		}

		ingredients := make(map[game.Item]int)
		for itemClass, amount := range raw.Ingredients {
			ingredients[items[itemClass]] = amount
		}

		products := make(map[game.Item]int)
		for itemClass, amount := range raw.Products {
			products[items[itemClass]] = amount
		}

		var power float64
		if raw.PowerConstant == 0 && raw.PowerFactor == 1 {
			power = float64(machine.Power)
		} else {
			power = float64(raw.PowerFactor)
		}

		recipeNames[raw.Name]++
		if recipeNames[raw.Name] > 1 {
			raw.Name = raw.Name + "(" + strconv.Itoa(recipeNames[raw.Name]) + ")"
		}

		recipes = append(recipes, &game.Recipe{
			Name:        raw.Name,
			Ingredients: ingredients,
			Products:    products,
			Duration:    float64(raw.Duration),
			Machine:     machine.Name,
			Power:       power,
		})
	}

	return recipes
}

func printRecipes(recipes []*game.Recipe) {
	recipeFile, err := os.Create("recipes.txt")
	if err != nil {
		panic(err)
	}
	defer recipeFile.Close()

	for _, recipe := range recipes {
		fmt.Fprintf(recipeFile, "%s\n\t%s %f %f\n\n", recipe.Name, recipe.Machine, recipe.Power, recipe.Duration)
		for product, amount := range recipe.Products {
			fmt.Fprintf(recipeFile, "\t%s: %d\n", product, amount)
		}
		fmt.Fprintln(recipeFile)
		for ingredient, amount := range recipe.Ingredients {
			fmt.Fprintf(recipeFile, "\t%s: %d\n", ingredient, amount)
		}
		fmt.Fprintln(recipeFile)
	}
}

type RawRecipe struct {
	ClassName     string        `json:"ClassName"`
	Name          string        `json:"mDisplayName"`
	Ingredients   ItemAmounts   `json:"mIngredients"`
	Products      ItemAmounts   `json:"mProduct"`
	Duration      StringFloat   `json:"mManufactoringDuration"`
	Machine       RecipeMachine `json:"mProducedIn"`
	PowerConstant StringFloat   `json:"mVariablePowerConsumptionConstant"`
	PowerFactor   StringFloat   `json:"mVariablePowerConsumptionFactor"`
}

type StringFloat float64

func (f *StringFloat) UnmarshalJSON(data []byte) error {
	var float float64
	err := json.Unmarshal(data[1:len(data)-1], &float)
	*f = StringFloat(float)
	return err
}

type RecipeMachine string

func (m *RecipeMachine) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 {
		return nil
	}
	data = data[:len(data)-2]
	*m = RecipeMachine(bytes.Split(bytes.Split(data, []byte(","))[0], []byte("."))[1])
	return nil
}

type ItemAmounts map[string]int

func (i *ItemAmounts) UnmarshalJSON(data []byte) error {
	if *i == nil {
		*i = make(ItemAmounts)
	}

	data = data[3 : len(data)-3]
	list := bytes.Split(data, []byte("),("))
	for _, raw := range list {
		split := bytes.Split(raw, []byte("\\\"',Amount="))
		amount, err := strconv.Atoi(string(split[1]))
		if err != nil {
			return err
		}

		if amount/1000 > 0 {
			amount /= 1000
		}

		split = bytes.Split(split[0], []byte("."))
		itemClass := string(split[len(split)-1])
		(*i)[itemClass] = amount
	}

	return nil
}
