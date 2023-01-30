package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	ads := getBuildInfo(AssemblyDirectorSystem, 4)
	mfg := getBuildInfo(MagneticFieldGenerator, 4)

	baseMaterials := make(map[Item]float64)
	getBaseMaterials(baseMaterials, ads.Tree)
	getBaseMaterials(baseMaterials, mfg.Tree)
	data, _ := json.MarshalIndent(baseMaterials, "", "\t")
	fmt.Println(string(data))
	fmt.Println()

	fmt.Println(ads)
	fmt.Println(mfg)
}

type BuildInfo struct {
	Tree *TreeNode
}

func (bi BuildInfo) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("PRODUCTION STEPS FOR %s %s/min\n\n", fmtAmount(bi.Tree.Output.Flux), bi.Tree.Output.Item))
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
	if isBase(node.Output.Item) {
		materials[node.Output.Item] += node.Output.Flux
		return
	}

	for _, input := range node.Inputs {
		getBaseMaterials(materials, input)
	}
}

func printBuildTree(buf *strings.Builder, node *TreeNode, indentation int) {
	prefix := '-'
	if isBase(node.Output.Item) {
		prefix = '■'
	}
	buf.WriteString(fmt.Sprintf("%s%c %s %s\n", strings.Repeat("\t", indentation), prefix, fmtAmount(node.Output.Flux), node.Output.Item))
	for _, child := range node.Inputs {
		printBuildTree(buf, child, indentation+1)
	}
}

type ItemFlux struct {
	Item Item
	Flux float64
}

type TreeNode struct {
	Output ItemFlux
	Inputs []*TreeNode
}

func getBuildInfo(item Item, amount float64) BuildInfo {
	info := BuildInfo{
		Tree: &TreeNode{
			Output: ItemFlux{
				item, amount,
			},
		},
	}

	_getRequirements(info.Tree)

	return info
}

func _getRequirements(node *TreeNode) {
	recipe, ok := allRecipes[node.Output.Item]
	if !ok {
		panic(node.Output.Item)
	}

	for inItem, inRate := range recipe.Input {
		totalInRate := node.Output.Flux * float64(inRate) / float64(recipe.Output[node.Output.Item])
		child := &TreeNode{
			Output: ItemFlux{inItem, totalInRate},
		}
		node.Inputs = append(node.Inputs, child)

		if !isBase(inItem) {
			_getRequirements(child)
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
