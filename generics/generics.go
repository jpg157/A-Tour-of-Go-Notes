package main

import "fmt"

// --- Type parameters in generic function or method ---

// Go functions can accept multiple types using type parameters
// The type parameters of a function appear between brackets before the function arguments

// This Index function works for any type that supports comparison (since type T implements comparable)
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// v and x are type T, which has the comparable
		// constraint, so we can use == here.
		if v == x {
			return i
		}
	}
	return -1
}

// The above declaration means that s is a slice of any type T that
// fulfills the built-in constraint comparable. x is also a value of the same type.

func main() {
	// index works on a slice of ints as well as slice of strings
	si := []int{10, 20, 15, -10}
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(si, 15))
	fmt.Println(Index(ss, "hi"))
}

// --- Generic Types ---

// In Go, a struct or interface can be parameterized with a type parameter,
// which can be useful for implementing generic data structures.

// List represents a singly-linked list that holds
// values of any type.
type LList[T comparable] struct {
	next *LList[T]
	val  T
}

type Flyer[T comparable] interface {
	Fly(distance T)
}

type PaginatedResDto[T any] struct {
	totalItems   int
	totalPages   int
	currPage     int
	itemsPerPage int
	nextPage     *int
	prevPage     *int
	data         []T
}
