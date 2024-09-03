package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for their name
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // Trim newline character

	// Prompt user for their age
	fmt.Print("Enter your age: ")
	ageInput, _ := reader.ReadString('\n')
	ageInput = strings.TrimSpace(ageInput) // Trim newline character

	age, err := strconv.Atoi(ageInput)
	if err != nil {
		fmt.Println("Invalid age. Please enter a number.")
		return
	}

	// Respond with a greeting
	fmt.Printf("Hello, %s! You are %d years old.\n", name, age)
}
