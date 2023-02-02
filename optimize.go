package main

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
