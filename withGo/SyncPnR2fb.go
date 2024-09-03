package main

import (
	"fmt"
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

// ... (previous type definitions remain the same)

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
	var lastCalculatedCount int

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
						lastCalculatedCount = len(fibSequence)
						return append(pnrs, 
							PnR{Name: "Average", Value: avg, Trivalent: "True"},
							PnR{Name: "LastCalculatedCount", Value: lastCalculatedCount, Trivalent: "True"})
					}
					return pnrs
				},
				Precondition: func(pnrs []PnR) bool {
					var fibSequence []int
					var fibRange []int
					for _, pnr := range pnrs {
						if pnr.Name == "FibSequence" {
							if seq, ok := pnr.Value.([]int); ok {
								fibSequence = seq
							}
						}
						if pnr.Name == "FibRange" {
							if rng, ok := pnr.Value.([]int); ok {
								fibRange = rng
							}
						}
					}
					
					// Check if we have a new Fibonacci number to process
					if len(fibSequence) > lastCalculatedCount {
						// Check if the last Fibonacci number is within the specified range
						if len(fibRange) == 2 && len(fibSequence) > 0 {
							lastFib := fibSequence[len(fibSequence)-1]
							return lastFib >= fibRange[0] && lastFib <= fibRange[1]
						}
						return true
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
            
            dc_activated := false
            
            // Check if we need to start a new IntentionLoop
            if len(cpux.IntentionLoop) == 0 {
                if cpux.DesignChunks[0].Precondition(sharedPnRs) {
                    cpux.IntentionLoop = []PnR{} // Initialize new IntentionLoop
                    dc_activated = true
                    cpux_activationFlag = true
                } else {
                    mutex.Unlock()
                    continue // Skip this CPUX if we can't start a new IntentionLoop
                }
            }

            // Process all DesignChunks if IntentionLoop is active
            for _, dc := range cpux.DesignChunks {
                if dc.Precondition(append(sharedPnRs, cpux.IntentionLoop...)) {
                    newPnRs := dc.Action(append(sharedPnRs, cpux.IntentionLoop...))
                    cpux.IntentionLoop = append(cpux.IntentionLoop, newPnRs...)
                    sharedPnRs = append(sharedPnRs, newPnRs...)
                    dc_activated = true
                }
            }

            // Clear IntentionLoop if no DesignChunks were activated this pass
            if !dc_activated {
                cpux.IntentionLoop = nil
            } else {
                cpux_activationFlag = true
            }

            mutex.Unlock()
        }
        
        if !cpux_activationFlag {
            break // Break the time-duration loop if no CPUX was activated
        }
        time.Sleep(100 * time.Millisecond) // Small delay to prevent tight loop
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