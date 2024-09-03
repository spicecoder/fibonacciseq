
package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
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
	Function func(*Robot, *sync.WaitGroup)
}

// Robot represents a robot runner (red or blue)
type Robot struct {
	Color          string
	Position       string
	BallsCollected int
	Speed          time.Duration
	NeedsRestart   bool
}

// Global PnR set
var globalPnR = map[string]*PnR{
	"BallsInArena":       {Name: "BallsInArena", Value: 20, Status: "True"},
	"RedRobotRunning":    {Name: "RedRobotRunning", Value: false, Status: "False"},
	"BlueRobotRunning":   {Name: "BlueRobotRunning", Value: false, Status: "False"},
	"RedRobotCollected":  {Name: "RedRobotCollected", Value: 0, Status: "True"},
	"BlueRobotCollected": {Name: "BlueRobotCollected", Value: 0, Status: "True"},
}

// Design Chunks
var designChunks = map[string]*DesignChunk{
	"Start": {
		Name: "Start",
		Function: func(r *Robot, wg *sync.WaitGroup) {
			r.Position = "Starting Point"
			globalPnR[r.Color+"RobotRunning"].Value = true
			globalPnR[r.Color+"RobotRunning"].Status = "True"
		},
	},
	"Run": {
		Name: "Run",
		Function: func(r *Robot, wg *sync.WaitGroup) {
			time.Sleep(r.Speed)
			r.Position = "Ball Collection Zone"
		},
	},
	"Collect": {
		Name: "Collect",
		Function: func(r *Robot, wg *sync.WaitGroup) {
			arenaPnR := globalPnR["BallsInArena"]
			if arenaPnR.Value.(int) > 0 {
				arenaPnR.Value = arenaPnR.Value.(int) - 1
				r.BallsCollected++
				globalPnR[r.Color+"RobotCollected"].Value = r.BallsCollected
			}
			time.Sleep(time.Millisecond * 500) // Time to collect the ball
		},
	},
	"Return": {
		Name: "Return",
		Function: func(r *Robot, wg *sync.WaitGroup) {
			time.Sleep(r.Speed)
			r.Position = "Starting Point"
		},
	},
}

func nameNorm(stringName string) string {
	stringName = strings.TrimSpace(stringName)
	spaceRegex := regexp.MustCompile(`\s+`)
	stringName = spaceRegex.ReplaceAllString(stringName, " ")
	alphanumRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	return alphanumRegex.ReplaceAllString(stringName, "")
}

func syncTest(gateMan, visitor map[string]*PnR) bool {
	for key, pnrA := range gateMan {
		normalizedKeyA := nameNorm(key)
		found := false
		for keyB, pnrB := range visitor {
			normalizedKeyB := nameNorm(keyB)
			if normalizedKeyA == normalizedKeyB {
				found = true
				trivalenceA := pnrA.Status
				if trivalenceA == "" {
					trivalenceA = "True"
				}
				trivalenceB := pnrB.Status
				if trivalenceB == "" {
					trivalenceB = "True"
				}
				if trivalenceA != trivalenceB {
					return false // Mismatch in trivalence portion
				}
				break
			}
		}
		if !found {
			return false // Corresponding PnR not found in visitor
		}
	}
	return true // All PnRs matched successfully
}

// IntentionLoop represents the execution of a CPUX
func IntentionLoop(robot *Robot, wg *sync.WaitGroup, restartChan chan bool, doneChan chan bool) {
	defer wg.Done()

	gateMan := map[string]*PnR{
		robot.Color + "RobotRunning": {Name: robot.Color + "RobotRunning", Value: true, Status: "True"},
		"BallsInArena":               {Name: "BallsInArena", Value: 20, Status: "True"},
	}

	for {
		if !syncTest(gateMan, globalPnR) {
			fmt.Printf("\n%s robot: PnR sync failed, waiting for alignment", robot.Color)
			time.Sleep(time.Second) // Wait before retrying
			continue
		}

		for _, chunk := range []string{"Start", "Run", "Collect", "Return"} {
			designChunks[chunk].Function(robot, wg)
		}

		gateMan[robot.Color+"RobotRunning"].Value = globalPnR[robot.Color+"RobotRunning"].Value
		gateMan[robot.Color+"RobotRunning"].Status = globalPnR[robot.Color+"RobotRunning"].Status
		gateMan["BallsInArena"].Value = globalPnR["BallsInArena"].Value
		gateMan["BallsInArena"].Status = globalPnR["BallsInArena"].Status

		if robot.BallsCollected >= 5 {
			robot.NeedsRestart = true
			globalPnR[robot.Color+"RobotRunning"].Value = false
			globalPnR[robot.Color+"RobotRunning"].Status = "False"
			fmt.Printf("\n%s robot needs restart after collecting 5 balls", robot.Color)

			select {
			case <-restartChan:
				fmt.Printf("\nSpace Loop restarting %s robot", robot.Color)
				robot.NeedsRestart = false
				robot.BallsCollected = 0
				globalPnR[robot.Color+"RobotRunning"].Value = true
				globalPnR[robot.Color+"RobotRunning"].Status = "True"
			case <-doneChan:
				fmt.Printf("\n%s robot shutting down", robot.Color)
				return
			}
		}

		if globalPnR["BallsInArena"].Value.(int) < 2 {
			fmt.Printf("\n%s robot stopped (less than 2 balls in arena)", robot.Color)
			return
		}
	}
}

// SpaceLoop coordinates the execution of all CPUX units
func SpaceLoop() {
	var wg sync.WaitGroup

	redRobot := &Robot{Color: "Red", Speed: time.Millisecond * time.Duration(rand.Intn(500) + 500)}
	blueRobot := &Robot{Color: "Blue", Speed: time.Millisecond * time.Duration(rand.Intn(500) + 500)}

	redRestartChan := make(chan bool)
	blueRestartChan := make(chan bool)
	redDoneChan := make(chan bool)
	blueDoneChan := make(chan bool)

	wg.Add(2)
	fmt.Println("Space Loop booting up Red robot")
	go IntentionLoop(redRobot, &wg, redRestartChan, redDoneChan)
	fmt.Println("Space Loop booting up Blue robot")
	go IntentionLoop(blueRobot, &wg, blueRestartChan, blueDoneChan)

	// Display and restart loop
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\rBalls in arena: %d | Red Robot: %d | Blue Robot: %d", 
				globalPnR["BallsInArena"].Value, 
				globalPnR["RedRobotCollected"].Value, 
				globalPnR["BlueRobotCollected"].Value)

			// Check and restart robots that need restarting
			if redRobot.NeedsRestart && globalPnR["BallsInArena"].Value.(int) >= 2 {
				fmt.Printf("\nSpace Loop detected Red robot needs restart. Initiating restart sequence...")
				redRestartChan <- true
			}
			if blueRobot.NeedsRestart && globalPnR["BallsInArena"].Value.(int) >= 2 {
				fmt.Printf("\nSpace Loop detected Blue robot needs restart. Initiating restart sequence...")
				blueRestartChan <- true
			}

			if globalPnR["BallsInArena"].Value.(int) < 2 {
				fmt.Println("\nLess than 2 balls in arena. Ending simulation...")
				close(redDoneChan)
				close(blueDoneChan)
				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Initializing Robot Sport Arena Simulation")
	fmt.Println("------------------------------------------")
	SpaceLoop()
	fmt.Println("------------------------------------------")
	fmt.Println("Simulation completed!")
}