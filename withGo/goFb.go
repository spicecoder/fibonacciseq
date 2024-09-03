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
	mutex             sync.Mutex
}

// Object structure (analogous to FBSequence in the JS example)
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
func (dc *DesignChunk) trigger(pnr *PnR, object *Object) {
	if dc.precondition(pnr) {
		fmt.Printf("%s precondition met. Executing action.\n", dc.name)
		dc.action(pnr, object)
	} else {
		fmt.Printf("%s precondition not met. Skipping action.\n", dc.name)
	}
}

// DesignChunk1 action to collect min and max values
func collectMinMax(pnr *PnR, object *Object) {
	fmt.Printf("DesignChunk1: Collecting min and max values.\n")

	if pnr.min == nil {
		var min int
		fmt.Print("Enter the minimum value: ")
		fmt.Scan(&min)
		pnr.min = &min
	}

	if pnr.max == nil {
		var max int
		fmt.Print("Enter the maximum value: ")
		fmt.Scan(&max)
		pnr.max = &max
	}

	// Emit the intention with the collected values
	payload := map[string]interface{}{
		"min": *pnr.min,
		"max": *pnr.max,
	}
	intention := &Intention{name: "setMinMax", payload: payload}
	object.receive(intention)
}

// DesignChunk2 action to generate Fibonacci sequence
func generateFibonacci(pnr *PnR, object *Object) {
	fmt.Printf("DesignChunk2: Generating Fibonacci sequence.\n")

	pnr.mutex.Lock()
	defer pnr.mutex.Unlock()

	x, y := 0, 1
	for x <= *pnr.max {
		if x >= *pnr.min {
			pnr.fibonacci = append(pnr.fibonacci, x)
		}
		x, y = y, x+y
	}

	// Set the maxIntReached condition to "yes" when the sequence is complete
	pnr.maxIntReached = "yes"
	pnr.isComplete = true

	// Print the Fibonacci sequence when complete
	fmt.Println("Fibonacci sequence in the range:", pnr.fibonacci)
}

// Precondition check for DesignChunk1
func preconditionMinMax(pnr *PnR) bool {
	return pnr.min == nil || pnr.max == nil
}

// Precondition check for DesignChunk2
func preconditionFibonacci(pnr *PnR) bool {
	return pnr.min != nil && pnr.max != nil && pnr.maxIntReached == "no"
}

// receive method for Object
func (o *Object) receive(intention *Intention) {
	fmt.Printf("Object received intention: %s\n", intention.name)

	o.pnr.mutex.Lock()
	defer o.pnr.mutex.Unlock()

	if intention.name == "setMinMax" {
		if min, ok := intention.payload["min"].(int); ok {
			o.pnr.min = &min
		}
		if max, ok := intention.payload["max"].(int); ok {
			o.pnr.max = &max
		}
	}
}

// Intention loop that triggers each design chunk
func intentionLoop(pnr *PnR, object *Object, designChunks []*DesignChunk) {
	for {
		executed := false

		for _, dc := range designChunks {
			dc.trigger(pnr, object)
			if dc.precondition(pnr) {
				executed = true
			}
		}

		if !executed {
			break // Stop the loop if no design chunk was executed
		}
	}
}

// main function to initialize and run the intention loop
func main() {
	pnr := PnR{
		maxIntReached: "no", // Initialize the condition to "no"
	}
	object := Object{pnr: &pnr}

	dc1 := DesignChunk{
		name:         "DesignChunk1",
		precondition: preconditionMinMax,
		action:       collectMinMax,
	}

	dc2 := DesignChunk{
		name:         "DesignChunk2",
		precondition: preconditionFibonacci,
		action:       generateFibonacci,
	}

	designChunks := []*DesignChunk{&dc1, &dc2}

	intentionLoop(&pnr, &object, designChunks)
}
