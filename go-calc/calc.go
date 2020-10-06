package main

import "strconv"

// Add two numbers represented as string
func Add(a string, b string) int64 {
	x, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		panic(err)
	}

	y, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		panic(err)
	}

	return x + y
}
