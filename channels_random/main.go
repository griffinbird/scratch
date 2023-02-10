package main

import (
	"log"

	"github.com/griffinbird/scratch/channels_random/helpers"
)

const numpool = 10

func main() {
	intChan := make(chan int)
	defer close(intChan)

	go CalculateValue(intChan)

	num := <-intChan
	log.Println(num)
}

func CalculateValue(intChan chan int) {
	randomNumber := helpers.RandomNumber(numpool)
	intChan <- randomNumber
}
