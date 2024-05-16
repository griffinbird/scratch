package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	low := 1
	high := 100
	tries := 0
	previousGuess := 0
	guess := 0

	fmt.Println("Please think of a number between", low, "and", high)
	fmt.Println("Press ENTER when ready")
	scanner.Scan()

	for {
		// binary search strategy
		previousGuess = guess
		guess = (low + high) / 2
		fmt.Println("I guess the number is", guess)
		tries++
		fmt.Println("Is that:")
		fmt.Println("(a) too high?")
		fmt.Println("(b) too low?")
		fmt.Println("(c) correct?")
		scanner.Scan()
		response := scanner.Text()

		if previousGuess == guess {
			fmt.Println("You're lying")
			os.Exit(1)
		}
		if response == "a" {
			high = guess - 1
		} else if response == "b" {
			low = guess + 1
		} else if response == "c" {
			fmt.Println("I won!")
			fmt.Println("It took", tries, "tries ito guess the correct answer")
			break
		} else {
			fmt.Println("Invalid response, try again.")
		}
	}
}
