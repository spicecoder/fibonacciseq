package main

import (
	"fmt"
	"sync"
)

// Intention structure
type Intention struct {
	name    string
	payload map[string]interface{}
}

// PnR structure
type PnR struct {
	min               *int
	max               *int
	fibonacci         []int
	isComplete        bool
	maxIntReached     string
	averageGenerated  string
	average           float64
	mutex             sync.Mutex
	executionStates   map[string]string // Track execution state of each design chunk
}

// Object structure
type Object struct {
	pnr *PnR
}

// Design Chunk structure
type DesignChunk struct {
	name         string
	precondition func(*PnR) bool
	action       func(*PnR, *Object)
}

// Method to trigger the precondition check and action
func (dc *DesignChunk) trigger(pnr *PnR, object *Object) bool {
	execKey := dc.name + " in execution"

	// Lock to check and set execution state
	pnr.mutex.Lock()
	if pnr.executionStates[execKey] == "Y" {
		fmt.Printf("%s is already in execution. Skipping.\n", dc.name)
		pnr.mutex.Unlock()
		return false
	}

	// Check if the precondition is met
	if dc.precondition(pnr) {
		fmt.Printf("PnR state at precondition check: %+v\n", *pnr)
		fmt.Printf("%s precondition met. Executing action.\n", dc.name)

		// Set the execution state to "Y" before executing the action
		pnr.executionStates[execKey] = "Y"
		pnr.mutex.Unlock() // Unlock before running the action

		// Execute the action
		dc.action(pnr, object)

		// Lock again to update execution state
		pnr.mutex.Lock()
		pnr.executionStates[execKey] = "N"
		fmt.Printf("PnR state after execution: %+v\n", *pnr)
		pnr.mutex.Unlock()

		return true
	} else {
		fmt.Printf("PnR state at precondition check: %+v\n", *pnr)
		fmt.Printf("%s precondition not met. Skipping action.\n", dc.name)
		pnr.mutex.Unlock()
		return false
	}
}

// CPUX structure
type CPUX struct {
	designChunks []*DesignChunk
	object       *Object
	pnr          *PnR
}

// DesignChunk1 of CPUX1: Collect min and max values
func collectMinMax(pnr *PnR, object *Object) {
	fmt.Printf("CPU1: Collecting min and max values.\n")

	if pnr.min == nil {
		var min int
		fmt.Print("Enter the minimum value: ")
		fmt.Scan(&min)
		pnr.mutex.Lock()
		pnr.min = &min
		pnr.mutex.Unlock()
	}

	if pnr.max == nil {
		var max int
		fmt.Print("Enter the maximum value: ")
		fmt.Scan(&max)
		pnr.mutex.Lock()
		pnr.max = &max
		pnr.mutex.Unlock()
	}

	// Emit the intention with the collected values
	payload := map[string]interface{}{
		"min": *pnr.min,
		"max": *pnr.max,
	}
	intention := &Intention{name: "setMinMax", payload: payload}
	object.receive(intention)
}

// DesignChunk2 of CPUX1: Generate Fibonacci sequence
func generateFibonacci(pnr *PnR, object *Object) {
	fmt.Printf("CPU1: Generating Fibonacci sequence.\n")

	x, y := 0, 1
	for x <= *pnr.max {
		if x >= *pnr.min {
			pnr.mutex.Lock()
			pnr.fibonacci = append(pnr.fibonacci, x)
			pnr.mutex.Unlock()
		}
		x, y = y, x+y
	}

	// Set the maxIntReached condition to "yes" when the sequence is complete
	pnr.mutex.Lock()
	pnr.maxIntReached = "yes"
	pnr.isComplete = true
	pnr.mutex.Unlock()

	// Print the Fibonacci sequence when complete
	fmt.Println("Fibonacci sequence in the range:", pnr.fibonacci)
}

// DesignChunk1 of CPUX2: Calculate the average of the Fibonacci numbers
func calculateAverage(pnr *PnR, object *Object) {
	fmt.Printf("CPU2: Calculating average of Fibonacci numbers.\n")

	pnr.mutex.Lock()
	defer pnr.mutex.Unlock()

	if len(pnr.fibonacci) == 0 {
		fmt.Println("No Fibonacci numbers available to calculate average.")
		return
	}

	sum := 0
	for _, num := range pnr.fibonacci {
		sum += num
	}

	pnr.average = float64(sum) / float64(len(pnr.fibonacci))
	pnr.averageGenerated = "yes"

	// Print the average
	fmt.Printf("Average of Fibonacci sequence: %.2f\n", pnr.average)
}

// Precondition check for DesignChunk1 of CPUX1
func preconditionMinMax(pnr *PnR) bool {
	execKey := "DesignChunk1 in execution"
	return (pnr.min == nil || pnr.max == nil) && pnr.executionStates[execKey] == "N"
}

// Precondition check for DesignChunk2 of CPUX1
func preconditionFibonacci(pnr *PnR) bool {
	execKey := "DesignChunk2 in execution"
	return pnr.min != nil && pnr.max != nil && pnr.maxIntReached == "no" && pnr.executionStates[execKey] == "N"
}

// Precondition check for DesignChunk1 of CPUX2
func preconditionAverage(pnr *PnR) bool {
	execKey := "DesignChunk1 in execution"
	return pnr.maxIntReached == "yes" && pnr.averageGenerated == "no" && pnr.executionStates[execKey] == "N"
}

// receive method for Object
func (o *Object) receive(intention *Intention) {
	o.pnr.mutex.Lock()
	defer o.pnr.mutex.Unlock()

	fmt.Printf("Object received intention: %s\n", intention.name)

	if intention.name == "setMinMax" {
		if min, ok := intention.payload["min"].(int); ok {
			o.pnr.min = &min
		}
		if max, ok := intention.payload["max"].(int); ok {
			o.pnr.max = &max
		}
	}
}

// Intention loop that triggers each design chunk in a CPUX
func intentionLoop(cpu *CPUX, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		executed := false

		for _, dc := range cpu.designChunks {
			if dc.trigger(cpu.pnr, cpu.object) {
				executed = true
			}
		}

		if !executed {
			fmt.Printf("Intention loop for %s stopped.\n", cpu.designChunks[0].name)
			return // Stop the loop if no design chunk was executed
		}
	}
}

// Space loop that triggers the first design chunk in each CPUX
func spaceLoop(cpuxs []*CPUX, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		executed := false

		for _, cpu := range cpuxs {
			// Trigger the first design chunk in each CPUX if precondition is met
			if len(cpu.designChunks) > 0 {
				if cpu.designChunks[0].trigger(cpu.pnr, cpu.object) {
					executed = true
				}
			}
		}

		if !executed {
			fmt.Println("Space loop stopped.")
			return // Stop the loop if no CPUX was triggered
		}
	}
}

func main() {
	// Initialize PnR and Objects
	pnr := PnR{
		maxIntReached:    "no", // Initialize the condition to "no"
		averageGenerated: "no", // Initialize the condition to "no"
		executionStates:  make(map[string]string),
	}

	// Initialize the executionStates for each DesignChunk
	pnr.executionStates["DesignChunk1 in execution"] = "N"
	pnr.executionStates["DesignChunk2 in execution"] = "N"

	object := Object{pnr: &pnr}

	// CPUX1: Fibonacci sequence generator
	cpu1 := CPUX{
		pnr:    &pnr,
		object: &object,
		designChunks: []*DesignChunk{
			{
				name:         "DesignChunk1",
				precondition: preconditionMinMax,
				action:       collectMinMax,
			},
			{
				name:         "DesignChunk2",
				precondition: preconditionFibonacci,
				action:       generateFibonacci,
			},
		},
	}

	// CPUX2: Average calculator for the Fibonacci sequence
	cpu2 := CPUX{
		pnr:    &pnr,
		object: &object,
		designChunks: []*DesignChunk{
			{
				name:         "DesignChunk1",
				precondition: preconditionAverage,
				action:       calculateAverage,
			},
		},
	}

	// Collection of CPUX
	cpuxs := []*CPUX{&cpu1, &cpu2}

	// Create a WaitGroup to wait for all loops to finish
	var wg sync.WaitGroup
	wg.Add(1) // Only add the space loop, which will trigger the intention loops

	// Run the space loop in a separate goroutine
	go spaceLoop(cpuxs, &wg)

	// Wait for all loops to complete
	wg.Wait()

	// Indicate that the program has completed successfully
	fmt.Println("All loops have completed. Program exiting.")
}
