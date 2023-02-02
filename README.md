# Overview

We've created a recipe analyzer for the game "Satisfactory". It takes in an item name and the desired amount to be produced, and outputs the required building and power resources as well as the recipes necessary to produce the desired item.

## Prerequisites

* Go 1.16 or later
* Satisfactory the Game
* The desire to built cool stuff

## Dependencies
*  github.com/PotatoesFall/satisfactory/game
*  github.com/erikgeiser/promptkit/selection
*  github.com/erikgeiser/promptkit/textinput
*  golang.org/x/exp/slices

 
## Usage

`go run main.go [item_name] [amount]`

_If_ no item and amount are provided, the code will prompt you to pick the item and enter the amount through the CLI interface.


## Current Limitations

Currently, the code does not include energy costs for resource extraction, building cost for resource extraction, or byproducts.

## Upcoming features

* Energy costs for resource extraction
* Building cost for resource extraction
* Accounting for byproducts 

