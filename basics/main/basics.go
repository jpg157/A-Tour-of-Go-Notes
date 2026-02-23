package main

// ^ All related source files in the same package need to have a
// package declaration statement at the top of the function

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const fileUploadSuccessfulMsg string = "Successfully created"
const fileUploadErrorMsg string = "An error occurred while attempting to create"

// Ex. function

func add(x, y int) int { // shorten all but the last func param if all are same
	return x + y
}

// Ex. returning multiple results from a function
// 		- return values may be named. If so, they are treated as variables defined at the top of the function (DO NOT NEED TO REDECLARE)

func swap(x, y string) (string, string) {
	return y, x
}

// Named return values

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x

	// A return statement without arguments returns the named return values (called "naked" return).
	// DO NOT USE THIS IN LONG FUNCTIONS, ONLY SHORT ONES (readability purposes)
	return
}

// Ex. var statement at package level (like global variables in C)
//
//	Don't do this. Should instead put mutable global variables in main.go and use dependency injection
var t1 string
var t2 int
var t3 bool

// Declaring a list of global variables of the SAME type in the same line
var c, python, java bool

// Ex. Multiple named return values and non-"naked" return
func test() (c bool, python bool, java string) {
	c, python, java = true, false, "no!"
	return c, python, java
}

// Ex. Can omit the type for a variable as the initializer is present

var variable_no_type = true

// := construct is not available outside of a function,
// as every statement outside should begin with a keyword (var, func, and so on)

// Ex. constant and variable declarations can be factored into "blocks"
// just like import statements
const (
	HelloConst1 = "hello" // Capital camel case for constants
	HelloConst2 = "hello"
	// (Constants can be character, string, boolean, or numeric values.
	// They can be inside or outside a function at the package level
	// Constants cannot be declared using the := syntax.)

	// An untyped constant takes the type needed by its context.
	Big   = 1 << 10  // should be a float64
	Small = Big >> 9 // should be an int since the value is 10 in binary or 2 in base 10
)

var (
	helloVar1 = "hello"
	helloVar2 = "hello"
)

// Unlike in C or Java, in Go assignment between items of different type requires an explicit conversion (no implicit widening or narrowing)
var x, y int = 3, 4
var f float64 = math.Sqrt(float64(x*x + y*y)) // need to explicitly cast
var z uint = uint(f)

// while in Go is a for loop without using the init and post condition components.
// No parenthesis around for loops in Go
// Ex.
func whileExample() {
	// For loop
	for i := 0; i < 10; i++ {
		fmt.Printf("hello %d", i)
	}

	// While loop

	var j int = 1
	for j < 10 {
		fmt.Printf("%d", j)
		j++
	}

	flag := false
	for flag != true {
		if 1 == (10 / 10) {
			flag = true
		}
	}
}

// Ex. a switch statement with no condition is the same as switch true
func switchExample() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

// Ex. Array syntax
func arrExample() {
	var a [2]string // declares an array of 2 strings.

	// Arrays cannot be resized in Go (size must be a constant known at compile time),
	// and are a contiguous memory block on the stack when allocated

	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Printf("%p\n", &a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)
}

type User struct {
	UserId string
	Name   string
}

const userId1 = "1d02455e-f24c-4c26-90d2-f1073c686314"
const userId2 = "96aeb270-dd19-4274-a2fe-30415644864b"

// Ex. Slice sntax
func sliceExample() {
	var primes [6]int = [6]int{2, 3, 5, 7, 11, 13}

	// Slice is a dynamic flexible view into the elements of an original array

	/*
		A slice does not store any data, it just describes a section of an underlying array (points to the array).

		Changing the elements of a slice modifies the corresponding elements of its underlying array.

		Other slices that share the same underlying array will see those changes.
	*/

	// --- Ways to make a slice: ---

	// 1. Using an existing array
	var pSlice []int = primes[0:4]

	// 2. Array-like declaration without init. with no size specified
	var slice2 []int
	fmt.Println(slice2)

	// Array-like declaration with init (called slice literal) - creates the same array as above, then builds a slice that references it
	var slice2Initialized []int = []int{1, 2, 3}
	fmt.Println(slice2Initialized)

	fmt.Println(pSlice)

	// Slice bound defaults ex.

	// Can omit the high or low bounds to use defaults instead.
	// Default is 0 for low bound and the length of the slice for high bound

	//  For the array
	// var a [10]int

	// these slice expressions are equivalent:
	// a[0:10]
	// a[:10]
	// a[0:]
	// a[:]

	// zero value of a slice is nil
	// a nil slice has a length and capacity of 0 and has no underlying array

	// 3. With make (allocates a zeroed array and returns a slice that references that array)
	makeSlice := make([]int, 4, 8)         // len(makeSlice) = 4, cap(makeSlice) = 8
	makeSlice = makeSlice[:cap(makeSlice)] // len(makeSlice)=5, cap(makeSlice)=5
	makeSlice = makeSlice[1:]              // len(makeSlice)=4, cap(makeSlice)=4

	// --- Slices of slices ex. ---
	// Create a tic-tac-toe board
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}
	board[0][1] = "X"

	// --- Appending to a slice ---

	// The first parameter s of append is a slice of type T, and the rest are T values to append to the slice.

	// The resulting value of append is a slice containing all the elements of the original slice plus the provided values.

	// If the backing array of s is too small to fit all the given values a bigger array will be allocated. The returned slice will point to the newly allocated array.
	var slice3 []User // len=0, cap=0 []
	slice3 = append(slice3, User{userId1, "John Doe"})
}

// Ex. range form of for loop
// iterates over a slice or map
func rangeForLoopEx() {

	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

	// two values are returned from each iteration. First is index, second is copy of the element at that index
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// Ex. can skip the index or value by assigning to _
	var records []User = []User{
		{UserId: userId1, Name: "John Doe"},
		{UserId: userId2, Name: "Jack Eod"},
	}

	for _, value := range records {
		fmt.Println(value.Name)
	}

	// Can omit the value second variable to only include the index
	for i := range pow {
		fmt.Println(i)
	}
}

// Ex. Maps
func mapExample() {
	// Maps keys to values.
	// a zero value of a map is nil, which has no keys, nor can keys be added

	// --- Ex. Creation of single map entry ---

	// the make function returns a map of the given type, initialized
	var dictionary map[string]string
	dictionary = make(map[string]string)
	dictionary["apple"] = "round, edible fruit of an apple tree"

	fmt.Println(dictionary["apple"])

	// --- Ex. Map literal ---

	// Like struct literal, but the keys are required
	var userLookupTable map[string]User
	userLookupTable = map[string]User{
		userId1: {UserId: userId1, Name: "John Doe"}, // if the top-level type is just a type name, you can omit it from the elements of the literal
		userId2: {UserId: userId2, Name: "Jack Eod"},
	}

	fmt.Println(userLookupTable[userId2])
	fmt.Println(userLookupTable[userId1])

	fmt.Println("map contents", userLookupTable)

	// --- Ex. map operations ---

	const newKey string = "orange"

	// Insert or update an element in map m
	dictionary[newKey] = "a round juicy citrus fruit with a tough bright reddish-yellow rind"

	// Retrieve an element
	var orangeDefinition string = dictionary[newKey]
	fmt.Println("orange dfn:", orangeDefinition)

	// Delete an element
	delete(dictionary, newKey)

	// Test that key is present with two-value assignment
	// If key is in the map, ok is true. If not, ok is false.
	// If key is not in the map, then elem is the zero value for the map's element type.
	elem, ok := dictionary[newKey]
	fmt.Println("The value:", elem, "Present?", ok)

	elem, ok = dictionary["apple"]
	fmt.Println("The value:", elem, "Present?", ok)
}

// Ex. Functions as values in Go

// In Go, functions are values too (first-class citizens).
// - Bc of this, they can be passed around just like other values.
// - They can be used as function arguments and return values
func FunctionValuesEx(fn func(x, y string) (string, string)) {
	str1 := "world"
	str2 := "hello"

	str1, str2 = fn(str1, str2)
	fmt.Println("Result after swap:", strings.Join(
		[]string{str1, str2},
		" "),
	)
}

// Ex. public and private functions, file upload

func HandleFileUpload(file string) string {
	var resMes string
	var uploadSuccess bool = storeFileInDb(file)

	if !uploadSuccess {
		resMes = fileUploadErrorMsg
	} else {
		resMes = fileUploadSuccessfulMsg
	}
	return resMes
}

func storeFileInDb(file string) bool {
	if file == "bad_file" {
		// fmt.Println("Log: Error while attempting to store file")
		return false
	}
	return true
}

func main() {
	// fmt.Println(add(42, 13))

	// // inside of a function, the := short assignment can be used in place
	// // of a variable with an implicit type (defined by the initializer)
	// res1, res2 := swap("1", "2")

	// fmt.Printf("Order after swapping: %s, %s\n", res1, res2)

	// var file string = "bad_file"
	// var fileUploadMessage string = HandleFileUpload(file)
	// fmt.Println(fileUploadMessage)

	// fmt.Println("hello", res1, "hi")

	// rangeForLoopEx()

	// mapExample()

	// Pass in the swap function as function argument
	FunctionValuesEx(swap)
}
