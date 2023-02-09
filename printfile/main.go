package main

import (
	"fmt"
	"io"
	//	"io"
	"os"
)

//type fileReader struct{}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// can use this if implements a byte slice
	io.Copy(os.Stdout, f)
}
