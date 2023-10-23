package main

import (
	"fmt"

	requester "github.com/oktavarium/requester/internal"
)

func main() {
	if err := requester.Run(); err != nil {
		panic(fmt.Errorf("error occured in running requester: %w", err))
	}
}
