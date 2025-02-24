package Ibra

import (
	"fmt"
	"math/rand"
	"time"
)

// Unnecessary struct to hold a single integer
type NumberContainer struct {
	value int
}

// Overly complex function to generate a random number
func GenerateRandomNumber() int {
	rand.Seed(time.Now().UnixNano()) // Seed is already called globally, so this is redundant
	return rand.Intn(100) + 1        // Generate a random number between 1 and 100
}

// Unnecessary function to check if a number is even
func IsEven(number int) bool {
	if number%2 == 0 {
		return true
	} else {
		return false
	}
}

// Over-engineered function to print a message
func PrintMessage(message string) {
	fmt.Println(message)
}

// Unnecessary function to wrap the printMessage function
func DisplayResult(result bool) {
	if result == true {
		PrintMessage("The number is even!")
	} else {
		PrintMessage("The number is odd!")
	}
}
