package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	_ "embed"
)

/*
TODO
- make the system be able to check different recipe variants
	- add weights for
		- resource rarity (per resource)
		- energy cost (don't forget miners/extractors)
		- number of machines? (per machine) nice-to-haven
- have the system account for by-products
	- match by-product recipes and balance them perfectly (how?)
	- nice-to-have: account for energy that can be produced with byproducts and deduct from total energy cost
- allow choosing/locking certain recipes while letting others be optimized?
*/

//go:embed recipes.json
var recipesJSON []byte

func main() {
	loadRecipes()

	ads := getBuildInfo(ItemFlux{Item: AssemblyDirectorSystem, Flux: 4})

	fmt.Println(ads)
}

type BuildInfo struct {
	Tree *TreeNode
}

func (bi BuildInfo) String() string {
	var buf strings.Builder
	// buf.WriteString(fmt.Sprintf("PRODUCTION STEPS FOR %s %s/min\n\n", fmtAmount(bi.Tree.Output.Flux), bi.Tree.Output.Item))
	baseMaterials := make(map[Item]float64)
	getBaseMaterials(baseMaterials, bi.Tree)
	buf.WriteString("Base Materials:\n")
	var materials []Item
	for mat := range baseMaterials {
		materials = append(materials, mat)
	}
	sort.Slice(materials, func(i, j int) bool { return materials[i] < materials[j] })
	for _, mat := range materials {
		buf.WriteString(fmt.Sprintf("- %s %s\n", fmtAmount(baseMaterials[mat]), mat))
	}
	buf.WriteString("\nProduction Tree:\n")
	printBuildTree(&buf, bi.Tree, 0)
	return buf.String()
}

func getBaseMaterials(materials map[Item]float64, node *TreeNode) {
	for _, output := range node.Output {
		if isBase(output.Item) {
			materials[output.Item] += output.Flux
			return
		}

		for _, input := range node.Input {
			getBaseMaterials(materials, input)
		}
	}
}

func printBuildTree(buf *strings.Builder, node *TreeNode, indentation int) {
	for _, output := range node.Output {
		prefix := '-'
		if isBase(output.Item) {
			prefix = '■'
		}
		buf.WriteString(fmt.Sprintf("%s%c %s %s\n", strings.Repeat("\t", indentation), prefix, fmtAmount(output.Flux), output.Item))
		for _, child := range node.Input {
			printBuildTree(buf, child, indentation+1)
		}
	}
}

type ItemFlux struct {
	Item Item
	Flux float64
}

type TreeNode struct {
	Output []ItemFlux
	Input  []*TreeNode
}

func getBuildInfo(items ...ItemFlux) BuildInfo {
	info := BuildInfo{
		Tree: &TreeNode{
			Output: items,
		},
	}

	_getRequirements(info.Tree)

	return info
}

func _getRequirements(node *TreeNode) {
	for _, output := range node.Output {
		recipes := recipesByItem[output.Item]
		if len(recipes) == 0 {
			panic("No recipes found for " + output.Item)
		}

		recipe := recipes[0] // TODO

		for inItem, inRate := range recipe.Input {
			totalInRate := output.Flux * float64(inRate) / float64(recipe.Output[output.Item])
			child := &TreeNode{
				Output: []ItemFlux{{inItem, totalInRate}},
			}
			node.Input = append(node.Input, child)

			if !isBase(inItem) {
				_getRequirements(child)
			}
		}
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
