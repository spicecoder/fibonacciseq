package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PnR represents a Prompt and Response pair
type PnR struct {
	Name   string
	Value  interface{}
	Status string // True/False/Undecided
}

// DesignChunk represents a unit of computation
type DesignChunk struct {
	Name     string
	Function func(*Runner, *sync.WaitGroup)
}

// Runner represents a runner (red or blue)
type Runner struct {
	Color       string
	Position    string
	BallsCollected int
	Speed       time.Duration
}

// Global PnR set
var globalPnR = map[string]*PnR{
	"BallsInBasket":      {Name: "BallsInBasket", Value: 20, Status: "True"},
	"RedRunnerRunning":   {Name: "RedRunnerRunning", Value: false, Status: "False"},
	"BlueRunnerRunning":  {Name: "BlueRunnerRunning", Value: false, Status: "False"},
	"RedRunnerCollected": {Name: "RedRunnerCollected", Value: 0, Status: "True"},
	"BlueRunnerCollected": {Name: "BlueRunnerCollected", Value: 0, Status: "True"},
}

// Design Chunks
var designChunks = map[string]*DesignChunk{
	"Start": {
		Name: "Start",
		Function: func(r *Runner, wg *sync.WaitGroup) {
			r.Position = "Starting Point"
			globalPnR[r.Color+"RunnerRunning"].Value = true
			globalPnR[r.Color+"RunnerRunning"].Status = "True"
		},
	},
	"Run": {
		Name: "Run",
		Function: func(r *Runner, wg *sync.WaitGroup) {
			time.Sleep(r.Speed)
			r.Position = "Basket"
		},
	},
	"Collect": {
		Name: "Collect",
		Function: func(r *Runner, wg *sync.WaitGroup) {
			basketPnR := globalPnR["BallsInBasket"]
			if basketPnR.Value.(int) > 0 {
				basketPnR.Value = basketPnR.Value.(int) - 1
				r.BallsCollected++
				globalPnR[r.Color+"RunnerCollected"].Value = r.BallsCollected
			}
			time.Sleep(time.Millisecond * 500) // Time to collect the ball
		},
	},
	"Return": {
		Name: "Return",
		Function: func(r *Runner, wg *sync.WaitGroup) {
			time.Sleep(r.Speed)
			r.Position = "Starting Point"
		},
	},
}

// IntentionLoop represents the execution of a CPUX
func IntentionLoop(runner *Runner, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		for _, chunk := range []string{"Start", "Run", "Collect", "Return"} {
			designChunks[chunk].Function(runner, wg)
		}

		if runner.BallsCollected >= 5 && globalPnR["BallsInBasket"].Value.(int) < 2 {
			globalPnR[runner.Color+"RunnerRunning"].Value = false
			globalPnR[runner.Color+"RunnerRunning"].Status = "False"
			break
		}
	}
}

// SpaceLoop coordinates the execution of all CPUX units
func SpaceLoop() {
	var wg sync.WaitGroup

	redRunner := &Runner{Color: "Red", Speed: time.Millisecond * time.Duration(rand.Intn(500) + 500)}
	blueRunner := &Runner{Color: "Blue", Speed: time.Millisecond * time.Duration(rand.Intn(500) + 500)}

	wg.Add(2)
	go IntentionLoop(redRunner, &wg)
	go IntentionLoop(blueRunner, &wg)

	// Display loop
	for {
		fmt.Printf("\rBalls in basket: %d | Red Runner: %d | Blue Runner: %d", 
			globalPnR["BallsInBasket"].Value, 
			globalPnR["RedRunnerCollected"].Value, 
			globalPnR["BlueRunnerCollected"].Value)

		if !globalPnR["RedRunnerRunning"].Value.(bool) && !globalPnR["BlueRunnerRunning"].Value.(bool) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
	fmt.Println("\nSimulation completed!")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	SpaceLoop()
}