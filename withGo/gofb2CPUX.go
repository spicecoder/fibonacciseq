package main

import (
	"fmt"
	"sync"
	"time"
)

// Intention structure
type Intention struct {
	name    string
	payload map[string]interface{}
}

// PnR structure
type PnR struct {
	mutex               sync.Mutex
	preconditions       map[string]string
	executionStates     map[string]string
	intentionloopActive string
	min, max            int
	fibonacci           []int
}

// Object structure
type Object struct {
	pnr *PnR
}

// DesignChunk structure
type DesignChunk struct {
	name         string
	precondition func(*PnR) bool
	action       func(*PnR, *Object)
}

// Method to trigger the precondition check and action asynchronously
func (dc *DesignChunk) trigger(pnr *PnR, object *Object, wg *sync.WaitGroup) {
	defer wg.Done()
	execKey := dc.name + " in execution"

	pnr.mutex.Lock()
	if pnr.executionStates[execKey] == "Y" {
		fmt.Printf("%s is already in execution. Skipping.\n", dc.name)
		pnr.mutex.Unlock()
		return
	}

	// Check if the precondition is met
	if dc.precondition(pnr) {
		fmt.Printf("%s precondition met. Executing action.\n", dc.name)

		// Set the execution state to "Y" and release the lock
		pnr.executionStates[execKey] = "Y"
		pnr.mutex.Unlock()

		// Start the asynchronous action
		dc.action(pnr, object)

		// After the action is complete, update the state
		pnr.mutex.Lock()
		pnr.executionStates[execKey] = "N"
		pnr.mutex.Unlock()
	} else {
		fmt.Printf("%s precondition not met. Skipping action.\n", dc.name)
		pnr.mutex.Unlock()
	}
}

// CPUX structure
type CPUX struct {
	designChunks []*DesignChunk
	object       *Object
	pnr          *PnR
}

// Start function to initiate the IntentionLoop
func (cpu *CPUX) start(wg *sync.WaitGroup) {
	defer wg.Done()

	cpu.pnr.mutex.Lock()
	cpu.pnr.intentionloopActive = "Y"
	cpu.pnr.mutex.Unlock()

	// Start the IntentionLoop
	intentionLoop(cpu)
}

// Intention loop that triggers each design chunk in a CPUX
func intentionLoop(cpu *CPUX) {
	for {
		executed := false

		// Visit each DesignChunk in the CPUX
		for _, dc := range cpu.designChunks {
			var wg sync.WaitGroup
			wg.Add(1)
			go dc.trigger(cpu.pnr, cpu.object, &wg)
			wg.Wait()

			// If any design chunk executes, set executed to true
			cpu.pnr.mutex.Lock()
			if cpu.pnr.executionStates[dc.name+" in execution"] == "Y" {
				executed = true
			}
			cpu.pnr.mutex.Unlock()

			// Configurable long wait after each design chunk
			time.Sleep(10 * time.Second)
		}

		if !executed {
			fmt.Println("No DesignChunks executed. Stopping IntentionLoop.")
			cpu.pnr.mutex.Lock()
			cpu.pnr.intentionloopActive = "N"
			cpu.pnr.mutex.Unlock()
			return
		}
	}
}

// Space loop that triggers the start function of each CPUX periodically
func spaceLoop(cpuxs []*CPUX, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		executed := false

		for _, cpu := range cpuxs {
			var wg sync.WaitGroup
			wg.Add(1)
			go cpu.start(&wg)
			wg.Wait()

			cpu.pnr.mutex.Lock()
			if cpu.pnr.intentionloopActive == "Y" {
				executed = true
			}
			cpu.pnr.mutex.Unlock()
		}

		// Wait between SpaceLoop iterations (configurable)
		time.Sleep(15 * time.Second)

		if !executed {
			fmt.Println("All CPUXs are in a non-executing state. Stopping SpaceLoop.")
			return
		}
	}
}

// receive method for Object
func (o *Object) receive(intention *Intention) {
	o.pnr.mutex.Lock()
	defer o.pnr.mutex.Unlock()

	fmt.Printf("Object received intention: %s\n", intention.name)

	// Handle intention and update PnR
	switch intention.name {
	case "setMinMax":
		if min, ok := intention.payload["min"].(int); ok {
			o.pnr.min = min
		}
		if max, ok := intention.payload["max"].(int); ok {
			o.pnr.max = max
		}
	}
}

// Main function to drive the program
func main() {
	// Initialize PnR and Objects
	pnr := PnR{
		preconditions:       make(map[string]string),
		executionStates:     make(map[string]string),
		intentionloopActive: "N",
	}

	// Initialize the executionStates and preconditions for each DesignChunk
	pnr.executionStates["DesignChunk1 in execution"] = "N"
	pnr.executionStates["DesignChunk2 in execution"] = "N"
	pnr.preconditions["DesignChunk1"] = "Y" // Start with DesignChunk1 ready to execute

	object := Object{pnr: &pnr}

	// DesignChunk1 asks for user input
	designChunk1 := &DesignChunk{
		name: "DesignChunk1",
		precondition: func(pnr *PnR) bool {
			return pnr.preconditions["DesignChunk1"] == "Y"
		},
		action: func(pnr *PnR, object *Object) {
			fmt.Print("Enter the minimum value: ")
			fmt.Scan(&pnr.min)
			fmt.Print("Enter the maximum value: ")
			fmt.Scan(&pnr.max)

			// Emit intention to set min and max values
			intention := &Intention{
				name: "setMinMax",
				payload: map[string]interface{}{
					"min": pnr.min,
					"max": pnr.max,
				},
			}
			object.receive(intention)

			// Update preconditions for the next DesignChunk
			pnr.preconditions["DesignChunk2"] = "Y"
			pnr.preconditions["DesignChunk1"] = "N"
		},
	}

	// DesignChunk2 generates the Fibonacci sequence
	designChunk2 := &DesignChunk{
		name: "DesignChunk2",
		precondition: func(pnr *PnR) bool {
			return pnr.preconditions["DesignChunk2"] == "Y"
		},
		action: func(pnr *PnR, object *Object) {
			fmt.Println("DesignChunk2 is executing...")
			// Generate Fibonacci sequence within the given range
			fibonacci := []int{}
			a, b := 0, 1
			for a <= pnr.max {
				if a >= pnr.min {
					fibonacci = append(fibonacci, a)
				}
				a, b = b, a+b
			}

			// Store the sequence in the PnR
			pnr.fibonacci = fibonacci
			fmt.Printf("Fibonacci sequence in range [%d, %d]: %v\n", pnr.min, pnr.max, pnr.fibonacci)

			// No further actions, so stop the IntentionLoop
			pnr.preconditions["DesignChunk2"] = "N"
		},
	}

	cpu := &CPUX{
		pnr:          &pnr,
		object:       &object,
		designChunks: []*DesignChunk{designChunk1, designChunk2},
	}

	// Collection of CPUX
	cpuxs := []*CPUX{cpu}

	// Create a WaitGroup to wait for the SpaceLoop to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the SpaceLoop in a separate goroutine
	go spaceLoop(cpuxs, &wg)

	// Wait for the SpaceLoop to complete
	wg.Wait()

	fmt.Println("All loops have completed. Program exiting.")
}
