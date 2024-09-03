package main

import (
	"fmt"
	"sync"
	"time"
)

// MemoryStore simulates a simple in-memory storage
type MemoryStore struct {
	mutex sync.Mutex
	name  string
	cond  *sync.Cond
}

// SetName safely sets the user's name in the memory store and notifies the greeting loop
func (store *MemoryStore) SetName(name string) {
	store.mutex.Lock()
	store.name = name
	store.mutex.Unlock()

	// Notify the greeting loop that a name has been set
	store.cond.Signal()
}

// GetName safely retrieves the user's name from the memory store
func (store *MemoryStore) GetName() string {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.name
}

// ClearName safely clears the user's name from the memory store
func (store *MemoryStore) ClearName() {
	store.mutex.Lock()
	store.name = ""
	store.mutex.Unlock()
}

// Function to ask the user for their name continuously
func askUserName(store *MemoryStore) {
	for {
		var name string
		fmt.Print("Enter your name: ")
		fmt.Scan(&name)
		if name != "" {
			store.SetName(name)

			// Wait for the greeting loop to process the current name before asking for a new one
			store.mutex.Lock()
			for store.name != "" {
				store.cond.Wait()
			}
			store.mutex.Unlock()
		}
	}
}

// Independent loop that greets the user periodically, but can be notified to act immediately
func greetingLoop(store *MemoryStore) {
	for {
		store.mutex.Lock()
		// Wait for the signal that a name has been entered
		for store.name == "" {
			store.cond.Wait()
		}

		name := store.name
		store.mutex.Unlock()

		// Check if the name is set
		if name != "" {
			fmt.Printf("Hello, %s!\n", name)

			// Delay before clearing the name to ensure the greeting is printed
			time.Sleep(1 * time.Second)
			store.ClearName()

			// Notify the askUserName function that the name has been processed
			store.cond.Signal()
			fmt.Println("Name cleared from store. Waiting for a new name...")
		} else {
			fmt.Println("Waiting for a name to greet...")
		}
	}
}

func main() {
	// Initialize the memory store
	store := &MemoryStore{}
	store.cond = sync.NewCond(&store.mutex)

	// Start the greeting loop in a separate goroutine
	go greetingLoop(store)

	// Continuously ask the user for their name
	askUserName(store)
}
