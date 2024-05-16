package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func main() {
	jim := person{
		firstName: "Jim",
		lastName:  "party",
		contactInfo: contactInfo{
			email:   "jim@gmail.com",
			zipCode: 9400,
		},
	}
	jim.updateName("Jimmy")
	jim.print()

}

func (p *person) updateName(newFirstName string) {

	p.firstName = newFirstName
}

func (p person) print() {
	fmt.Printf("Address of struct = %+v: %p\n", p, &p)
}