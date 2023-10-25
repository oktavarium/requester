package main

import (
	requester "github.com/oktavarium/requester/internal"
)

func main() {
	if err := requester.Run(); err != nil {
		panic(err)
	}
}
