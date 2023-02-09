package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	for l := range c {
		// funciton liternal / anonymous function
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, c)
		}(l) // passing in copy of l
	}
}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Print(link, "might be done!")
		c <- link
		return // so we don't do anything else
	}

	fmt.Println(link, "is up!")
	c <- link
}