package main

import (
	"fmt"
	"sync"
	"time"
)

type PnRValue struct {
	Answer     string
	Trivalence string
	Completed  bool
}

type PnR map[string]PnRValue

type DesignChunk struct {
	Name string
	PnR  PnR
}

type CPUX struct {
	Name          string
	DesignChunks  []DesignChunk
	CurrentChunk  int
	IsActive      bool
	IntentionLoop chan bool
}

func nameNorm(s string) string {
	return s // Simplified for this example
}

func syncTest(gateMan, visitor PnR) bool {
	for key, gateManValue := range gateMan {
		normalizedKeyA := nameNorm(key)
		found := false

		for visitorKey, visitorValue := range visitor {
			normalizedKeyB := nameNorm(visitorKey)

			if normalizedKeyA == normalizedKeyB {
				found = true
				if gateManValue.Trivalence != visitorValue.Trivalence || visitorValue.Completed {
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

func flowinPnR(dcPnR, rtPnR PnR) PnR {
	filteredPnR := make(PnR)
	for key, dcValue := range dcPnR {
		normalizedKey := nameNorm(key)
		for rtKey, rtValue := range rtPnR {
			normalizedRtKey := nameNorm(rtKey)
			if normalizedKey == normalizedRtKey && !rtValue.Completed {
				filteredPnR[key] = PnRValue{
					Answer:     dcValue.Answer,
					Trivalence: rtValue.Trivalence,
					Completed:  false,
				}
			}
		}
	}
	return filteredPnR
}

func flowoutPnR(dcPnR, rtPnR PnR) PnR {
	for key, dcValue := range dcPnR {
		normalizedKey := nameNorm(key)
		for rtKey, rtValue := range rtPnR {
			normalizedRtKey := nameNorm(rtKey)
			if normalizedKey == normalizedRtKey && !rtValue.Completed {
				rtPnR[rtKey] = PnRValue{
					Answer:     rtValue.Answer,
					Trivalence: dcValue.Trivalence,
					Completed:  true,
				}
				break
			}
		}
	}
	return rtPnR
}

func activityTest(items []interface{}, globalPnR PnR) bool {
	for _, item := range items {
		switch v := item.(type) {
		case *CPUX:
			if v.IsActive {
				return true
			}
		case DesignChunk:
			if syncTest(v.PnR, globalPnR) {
				return true
			}
		}
	}
	return false
}

func runIntentionLoop(cpux *CPUX, globalPnR PnR, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-cpux.IntentionLoop:
			return
		default:
			currentChunk := &cpux.DesignChunks[cpux.CurrentChunk]
			if syncTest(currentChunk.PnR, globalPnR) {
				fmt.Printf("Executing %s in CPUX %s\n", currentChunk.Name, cpux.Name)
				currentChunk.PnR = flowinPnR(currentChunk.PnR, globalPnR)
				time.Sleep(time.Millisecond * 100) // Simulating work
				globalPnR = flowoutPnR(currentChunk.PnR, globalPnR)
			}

			cpux.CurrentChunk = (cpux.CurrentChunk + 1) % len(cpux.DesignChunks)

			// Check if all chunks are completed
			allCompleted := true
			for _, chunk := range cpux.DesignChunks {
				for _, pnrValue := range chunk.PnR {
					if !pnrValue.Completed {
						allCompleted = false
						break
					}
				}
				if !allCompleted {
					break
				}
			}

			if allCompleted {
				cpux.IsActive = false
				return
			}

			time.Sleep(time.Millisecond * 10)
		}
	}
}

func runSpaceLoop(cpuxs []*CPUX, globalPnR PnR) {
	var wg sync.WaitGroup

	for _, cpux := range cpuxs {
		wg.Add(1)
		go runIntentionLoop(cpux, globalPnR, &wg)
	}

	for {
		if !activityTest(func() []interface{} {
			items := make([]interface{}, len(cpuxs))
			for i, cpux := range cpuxs {
				items[i] = cpux
			}
			return items
		}(), globalPnR) {
			fmt.Println("No active CPUXs, stopping Space Loop")
			for _, cpux := range cpuxs {
				close(cpux.IntentionLoop)
			}
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	wg.Wait()
	fmt.Println("Final Global PnR state:")
	for key, value := range globalPnR {
		fmt.Printf("%s: %+v\n", key, value)
	}
}

var globalPnR = PnR{
	"Question 1": PnRValue{Answer: "Answer 1", Trivalence: "True", Completed: false},
	"Question 2": PnRValue{Answer: "Answer 2", Trivalence: "False", Completed: false},
	"Question 3": PnRValue{Answer: "Answer 3", Trivalence: "Undecided", Completed: false},
}

func main() {
	cpux1 := &CPUX{
		Name: "CPUX1",
		DesignChunks: []DesignChunk{
			{Name: "DC1", PnR: PnR{"Question 1": PnRValue{Answer: "Answer 1", Trivalence: "True", Completed: false}}},
			{Name: "DC2", PnR: PnR{"Question 2": PnRValue{Answer: "Answer 2", Trivalence: "False", Completed: false}}},
		},
		IsActive:      true,
		IntentionLoop: make(chan bool),
	}

	cpux2 := &CPUX{
		Name: "CPUX2",
		DesignChunks: []DesignChunk{
			{Name: "DC3", PnR: PnR{"Question 3": PnRValue{Answer: "Answer 3", Trivalence: "Undecided", Completed: false}}},
			{Name: "DC4", PnR: PnR{"Question 1": PnRValue{Answer: "Answer 1", Trivalence: "True", Completed: false}}},
		},
		IsActive:      true,
		IntentionLoop: make(chan bool),
	}

	cpuxs := []*CPUX{cpux1, cpux2}

	runSpaceLoop(cpuxs, globalPnR)
}