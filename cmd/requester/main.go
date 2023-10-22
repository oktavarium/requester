package main

import "fmt"

func main() {
	if err := requester.Run(); err != nil {
		panic(fmt.Errorf("error occured in running requester: %w", err))
	}
}
