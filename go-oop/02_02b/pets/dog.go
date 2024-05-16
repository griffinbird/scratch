package pets

import (
	"fmt"
	"strings"
	"time"
)

type Dog struct {
	Name string
	Color string
	Breed string
	lastSlept time.Time
}

func (d Dog) needsSleep() bool {
	return time.Now().Sub(d.lastSlept) > 4 * time.Hour
}

func (d Dog) Sleep() {
	d.lastSlept = time.Now()
}

func (d Dog) Feed(food string) string {
	return fmt.Sprintf("%s is eating %s", d.Name, food)
}

func (d Dog) GiveAttention(activity string) string {
	if d.needsSleep() {
		d.Sleep()
		return "Your dog is tired and need to rest"
	}
	response := ""
	switch strings.ToUpper(activity) {
	case "PET":
		response = fmt.Sprintf("wags tail")
	case "Playing Fetch":
		response = "return the ball and jump waiting for you to throw it again"
	default:
		response = "bark"
	}
	return fmt.Sprintf("%s loves attention, %s will cause him to %s", d.Name, activity, response)
}

func NewDog (name, color, breed string, lastSlept time.Time) Dog {
	return Dog{
		Name: name,
		Color: color,
		Breed: breed,
		lastSlept: lastSlept,
	}
}