package main

import (
	"fmt"

	requester "github.com/oktavarium/requester/internal"
)

func main() {
	if err := requester.Run(); err != nil {
		fmt.Printf("error occured in running requester: %s", err)
	}
}
