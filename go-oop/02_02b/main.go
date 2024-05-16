package main

import (
	"fmt"
	"go-opp/pets"
	"time"
)

func main() {
	sleepTime := time.Now().Add(time.Duration(-5) * time.Hour)
	pet := pets.NewDog("Sunny", "Brown", "Golden Retriever", sleepTime)
	fmt.Println(pet.Feed("steak"))
	fmt.Println(pet.GiveAttention("play fetch"))
}
