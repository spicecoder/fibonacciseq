package main

import (
	"fmt"
	// "math"
	"strings"
	"sync"
	"time"
)

// PnR represents a Prompt and Response pair
type PnR struct {
	Name      string
	Value     interface{}
	Trivalent string // "True", "False", or "Undecided"
}

// DesignChunk represents a unit of computation
type DesignChunk struct {
	Name         string
	Action       func([]PnR) []PnR
	Precondition func([]PnR) bool
}

// CPUX represents a Computational Path of Understanding and Execution
type CPUX struct {
	Name          string
	DesignChunks  []DesignChunk
	IntentionLoop []PnR
}

// Synchronicity function to check if PnRs match
func syncTest(gateMan, visitor []PnR) bool {
	for _, pnrA := range gateMan {
		found := false
		for _, pnrB := range visitor {
			if strings.EqualFold(strings.TrimSpace(pnrA.Name), strings.TrimSpace(pnrB.Name)) {
				found = true
				trivalenceA := pnrA.Trivalent
				if trivalenceA == "" {
					trivalenceA = "True"
				}
				trivalenceB := pnrB.Trivalent
				if trivalenceB == "" {
					trivalenceB = "True"
				}
				if trivalenceA != trivalenceB {
					return false
				}
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func createFibonacciCPUX() CPUX {
	var fibRange []int
	var lastFib, secondLastFib int
	var fibSequence []int
	var delay time.Duration

	return CPUX{
		Name: "FibonacciGenerator",
		DesignChunks: []DesignChunk{
			{
				Name: "GetRange",
				Action: func(pnrs []PnR) []PnR {
					fibRange = []int{1, 100}
					delay = 500 * time.Millisecond
					fmt.Println("FibonacciGenerator: Range set to", fibRange)
					return append(pnrs, PnR{Name: "FibRange", Value: fibRange, Trivalent: "True"})
				},
				Precondition: func(pnrs []PnR) bool {
					return len(fibRange) == 0
				},
			},
			{
				Name: "GenerateFib",
				Action: func(pnrs []PnR) []PnR {
					time.Sleep(delay)
					if len(fibSequence) == 0 {
						lastFib = 1
						fibSequence = append(fibSequence, lastFib)
					} else if len(fibSequence) == 1 {
						secondLastFib = lastFib
						lastFib = 1
						fibSequence = append(fibSequence, lastFib)
					} else {
						newFib := lastFib + secondLastFib
						secondLastFib = lastFib
						lastFib = newFib
						fibSequence = append(fibSequence, lastFib)
					}
					fmt.Printf("FibonacciGenerator: Generated %d. Sequence: %v\n", lastFib, fibSequence)
					return append(pnrs, PnR{Name: "FibSequence", Value: fibSequence, Trivalent: "True"})
				},
				Precondition: func(pnrs []PnR) bool {
					if len(fibSequence) < 2 {
						return true
					}
					nextFib := fibSequence[len(fibSequence)-1] + fibSequence[len(fibSequence)-2]
					return nextFib <= fibRange[1]
				},
			},
		},
		IntentionLoop: []PnR{},
	}
}

func createAverageCPUX() CPUX {
	return CPUX{
		Name: "AverageCalculator",
		DesignChunks: []DesignChunk{
			{
				Name: "CalculateAverage",
				Action: func(pnrs []PnR) []PnR {
					var fibSequence []int
					for _, pnr := range pnrs {
						if pnr.Name == "FibSequence" {
							if seq, ok := pnr.Value.([]int); ok {
								fibSequence = seq
							}
						}
					}
					if len(fibSequence) > 0 {
						sum := 0
						for _, num := range fibSequence {
							sum += num
						}
						avg := float64(sum) / float64(len(fibSequence))
						fmt.Printf("AverageCalculator: Count: %d, Average: %.2f\n", len(fibSequence), avg)
						return append(pnrs, PnR{Name: "Average", Value: avg, Trivalent: "True"})
					}
					return pnrs
				},
				Precondition: func(pnrs []PnR) bool {
					for _, pnr := range pnrs {
						if pnr.Name == "Average" {
							if pnr.Trivalent == "True" {
								return false
							}
						}
						if pnr.Name == "FibSequence" {
							return true
						}
					}
					return false
				},
			},
		},
		IntentionLoop: []PnR{},
	}
}

func SpaceLoop(cpuxUnits []CPUX, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	
	sharedPnRs := []PnR{}
	var mutex sync.Mutex
	
	start := time.Now()
	for time.Since(start) < duration {
		cpux_activationFlag := false
		for _, cpux := range cpuxUnits {
			
			mutex.Lock()
			for _, dc := range cpux.DesignChunks {
				//dc_activationFlag := false
				if dc.Precondition(append(sharedPnRs, cpux.IntentionLoop...)) {
					//dc_activationFlag = true;
					newPnRs := dc.Action(append(sharedPnRs, cpux.IntentionLoop...))
					sharedPnRs = append(sharedPnRs, newPnRs...)
					cpux.IntentionLoop = append(cpux.IntentionLoop, newPnRs...)
				}
				//if (len(cpux.IntentionLoop) == 0) { break}
				//fmt.Printf("cpux Inentionloop: Generated :",  cpux.IntentionLoop)
				cpux_activationFlag = true
				// if !dc_activationFlag {
				//  	break
				// }
			}
			mutex.Unlock()
			
		}
		if !cpux_activationFlag {
			fmt.Printf("space loop exit")
			break
		}
		
		//time.Sleep(100 * time.Millisecond) // Small delay to prevent tight loop
	}
}

func main() {
	fibCPUX := createFibonacciCPUX()
	avgCPUX := createAverageCPUX()

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("Starting Space Loop...")
	go SpaceLoop([]CPUX{fibCPUX, avgCPUX}, 10*time.Second, &wg)

	wg.Wait()
	fmt.Println("Space Loop finished.")
}