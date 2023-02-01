package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PotatoesFall/satisfactory/game"
)

func parseItems(data ...[]byte) ([]game.Item, map[string]game.Item) {
	var items []game.Item
	itemClasses := make(map[string]game.Item)
	for _, data := range data {
		var rawItems []RawItem
		if err := json.Unmarshal(data, &rawItems); err != nil {
			panic(err)
		}

		for _, item := range rawItems {
			items = append(items, item.Process())
			itemClasses[item.Class] = item.Process()
		}

	}

	return items, itemClasses
}

type RawItem struct {
	Class string `json:"ClassName"`
	Name  string `json:"mDisplayName"`
}

func (r RawItem) Process() game.Item {
	return game.Item(r.Name)
}

func printItems(items []game.Item) {
	itemFile, err := os.Create("items.txt")
	if err != nil {
		panic(err)
	}
	defer itemFile.Close()
	for _, item := range items {
		fmt.Fprintln(itemFile, item)
	}
}
