package main

import (
	"bytes"
	"encoding/json"
)

type Recipe struct {
	Name    string       `json:"name"`
	Input   map[Item]int `json:"input"`
	Output  map[Item]int `json:"output"`
	Machine Machine      `json:"machine"`
	Power   float64      `json:"power,omitempty"`
	Time    int          `json:"time"` // seconds
}

type Machine string // Refinery, Constructor, etc.

const (
	Smelter      Machine = "Smelter"
	Constructor  Machine = "Constructor"
	Packager     Machine = "Packager"
	Assembler    Machine = "Assembler"
	Foundry      Machine = "Foundry"
	Refinery     Machine = "Refinery"
	Manufacturer Machine = "Manufacturer"
	Blender      Machine = "Blender"

	ParticleAccelerator Machine = "Particle Accelerator" // power needed varies
)

func (m Machine) Power(recipe Recipe) float64 {
	switch m {
	case Smelter:
		return 4
	case Constructor:
		return 4
	case Packager:
		return 10
	case Assembler:
		return 15
	case Foundry:
		return 16
	case Refinery:
		return 30
	case Manufacturer:
		return 55
	case Blender:
		return 75

	case ParticleAccelerator:
		return recipe.Power
	}

	panic("unknown machine " + m)
}

type Power float64 // in MW

func (p Power) MarshalJSON() ([]byte, error) {
	str := fmtAmount(float64(p))
	return []byte("\"" + str + " MW\""), nil
}

func (p *Power) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSuffix(data, []byte(" MW\""))
	data = bytes.TrimPrefix(data, []byte("\""))

	var float float64
	err := json.Unmarshal(data, &float)
	*p = Power(float)
	return err
}
