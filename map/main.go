package main

import (
	"fmt"
	"log"
)

type User struct {
	Firstname string
	Lastname  string
}

func main() {
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4bf745",
		"white": "ffffff",
	}

	myMap := make(map[string]User)

	me := User{
		Firstname: "Ben",
		Lastname:  "Griffin",
	}

	myMap["me"] = me

	log.Println(myMap["me"].Firstname)

	printMap(colors)

}

func printMap(c map[string]string) {
	for color, hex := range c {
		fmt.Println("Hex code for", color, "is", hex)

	}
}
