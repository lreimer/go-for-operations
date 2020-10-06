package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	result := Add("1", "1")
	if result != 2 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", result, 2)
	}
}
