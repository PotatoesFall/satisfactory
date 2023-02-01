package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RawMachine struct {
	Class string      `json:"ClassName"`
	Name  string      `json:"mDisplayName"`
	Power StringFloat `json:"mPowerConsumption"`
}

func parseMachines(data ...[]byte) ([]RawMachine, map[string]RawMachine) {
	var machines []RawMachine
	machineClasses := make(map[string]RawMachine)
	for _, data := range data {
		var rawMachines []RawMachine
		if err := json.Unmarshal(data, &rawMachines); err != nil {
			panic(err)
		}

		for _, machine := range rawMachines {
			machines = append(machines, machine)
			machineClasses[machine.Class] = machine
		}
	}

	return machines, machineClasses
}

func printMachines(machines []RawMachine) {
	machineFile, err := os.Create("machines.txt")
	if err != nil {
		panic(err)
	}
	defer machineFile.Close()
	for _, machine := range machines {
		fmt.Fprintln(machineFile, machine.Name, machine.Power)
	}
}
