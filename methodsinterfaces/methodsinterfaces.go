package main

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Vertex struct {
	X, Y float64
}

// === Methods ===

// Ex. Defining a method on a type
// - A method in Go is a function with a special reciever argument
// - The receiver appears in its own argument list between the func keyword and the method name.

// Ex. abs method has a receiver of type Vertex named v
func (v Vertex) abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Ex. Can declare a method on non-struct types too
type Latitude float64

func (f Latitude) toMetres() float64 {
	const MetresPerDegree int = 111132
	return float64(f * Latitude(MetresPerDegree))
}

// Note: Cannot declare a method with a receiver whose type is defined in another package
// (which includes the built-in types such as int).

// --- Pointer receivers ---

// Method with value reciever
func (v Vertex) scaleCopy(f float64) Vertex {
	v.X = v.X * f
	v.Y = v.Y * f
	return v
}

// vs

// Can declare methods with pointer recievers, in order to allow modifying the original value
// (rather than creating a copy of the original then only modifying that copy)
func (v *Vertex) scaleOriginal(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func methodExamples() {
	var v Vertex = Vertex{3, 4}

	copyOfV := v.scaleCopy(10)

	fmt.Printf("Original vertex val after scaling using value reciever method: %g\n", v.abs())
	fmt.Printf("Copy of vertex val after scaling using value reciever method: %g\n", copyOfV.abs())

	// --- Methods and pointer indirection (does not work for functions) ---

	// Ex. Methods with pointer receievers can accept values without having to explicitly create a reference
	v.scaleOriginal(10) // implicitly passes the address of v (returning a pointer to v) into the type method
	// "For the above statement, even though v is a value and not a pointer,
	// the method with the pointer receiver is called automatically.
	// Go interprets the statement v.ScaleOriginal(10) as (&v).ScaleOriginal(10)
	// since the ScaleOriginal method has a pointer receiver."

	// Ex. Works the other way too - value recievers can accept pointers to values without explicit dereference (excluding primitive types)
	var num Latitude = 5.0
	var pNum *Latitude = &num
	pNum.toMetres()
	// ^ works the same as:
	(*pNum).toMetres()

	// "In general, all methods on a given type should have either value or pointer receivers,
	// but not a mixture of both"

	fmt.Printf("Original vertex val after scaling using pointer reciever method: %g\n", v.abs())
}

//  === Interfaces ===

// In Go:
// An interface type defines a set of method signatures. It cannot have fields
// Interfaces are implemented implicitly
// A value of interface <type> can hold any value that implements all of the methods of that interface
// A single type can implement any number of interfaces as long as there are not method signature conflicts
// Naming convention;
// - For single method interfaces - use agent nouns, with -er at the end that describes functionality (ex. Reader, Flyer)
// - For multi-method interfaces - usename that describes overall purpose (noun or agent noun)
// Architecture convention is to define interfaces in the consumer package (or module if cross-module),
// not where the interface is implemented
type Abser interface {
	abs() float64
}

type I interface {
	PrintX()
}

func (v *Vertex) PrintX() {
	if v == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(v.X)
}

type Person struct {
	Name string
	Age  int
}

// --- The Stringer interface (Go's equivalent of toString() in Java or C#) ---
// - Defined in fmt package
// - Stringer is a type that can define itself as a string.
// - The fmt package uses this interface to print values (interface is defined where it is used)
// - You create the implementation for a concrete type
func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

func interfaceExamples() {
	var a Abser
	var v Vertex = Vertex{3, 4} // (side note) creating a Vertex using untyped numeric constants 3 and 4, which is implicitly converted to float64 (required in parameter list)

	// Vertex implicitly implements Abser (has abs() method defined on Vertex type - line 19),
	// so can be assigned to Abser (allows for polymorphism and following LSPrinciple)
	a = v
	fmt.Printf("(%v, %T)", a, a)
	fmt.Println(a.abs())

	// --- Interface values ---
	// Under the hood, interface values can be thought of as a tuple of a value and a concrete type:
	// (value, type)
	// An interface value holds a value of a specific underlying concrete type.

	// --- Interface values with nil underlying values ---
	// If concrete value in interface is nil, the method will be called with a nil receiver (rather than a null pointer exception occuring)
	var i I
	var v2 *Vertex
	i = v2
	i.PrintX()

	// same behaviour as above occurs when calling method on uninitialized struct
	v2.PrintX()

	// Ex. attempting to dereference zero value pointer of type Vertex, to access struct field
	// 		results in panic: runtime error: invalid memory address or nil pointer dereferences
	// fmt.Println(v2.X)

	// (Note - an interface value that holds a nil concrete value is itself non-nil)

	// --- Nil interface values ---

	// a nil interface holds neither interface value nor concrete type
	// calling a method on a nil interface results in a run-time error
	// (no interface tuple to indicate which concrete method to call)
	// var i2 I
	// i2.PrintX()

	// --- Empty interfaces ---

	// The interface type that specifies zero methods is known as the empty interface
	// An empty interface may hold values of any type. (Every type implements at least zero methods.)
	// Used by code that handles values of unknown type (ex. fmt.Print takes any number of args of type interface{}
	var anyTypeValue1 interface{}

	// in modern Go, any is used (alias of interface{})
	var anyTypeValue2 any = Vertex{}

	fmt.Printf("value of empty interface after no initialization: %v\n", anyTypeValue1)
	fmt.Printf("value empty interface after initialization to Vertex concrete type: %v\n", anyTypeValue2)

	// --- Type assertions ---

	// can access any interface value be making type assertion (like in Typescript - as string or <string>value)
	t2 := anyTypeValue2.(Vertex)

	// if the interface value does not hold a Vertex as concrete type, the statement will trigger a "panic",
	// unless a type assertion test syntax is used (with two values returned), to test if the interface holds the concrete type

	// If empty interface does not hold type,
	// then t will be zero value of that concrete type (no panic occurs) and ok as false
	// (similar to reading map value)
	t1, ok := anyTypeValue1.(Vertex)

	fmt.Println("Value of t2 after type assertion with no test", t2)
	fmt.Printf("Value of t1 after type assertion with test: %v | Type assertion passed: %t\n", t1, ok)

	// --- Type switches ---

	// Like regular switch statement, but cases in a type switch specify types (not values)
	// The declaration in a type switch has the same syntax as a type assertion i.(T),
	// but the specific type T is replaced with the keyword type.
	switch ctype := anyTypeValue2.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", ctype, ctype*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", ctype, len(ctype))
	default:
		fmt.Printf("Is non-primitive type. Type is %T and value is: %#v\n", ctype, ctype)
	}

	// --- Stringer interface continued ---

	var p Person = Person{"John Doe", 35}
	fmt.Printf("Custom string format of type %T:\n%v\n", p, p)

	// --- The error interface ---

	// Go programs express error state with error values.

	// The error type is a built-in interface similar to fmt.Stringer
	// (As with fmt.Stringer, the fmt package looks for the error interface when printing values.)
	// Functions often return an error value, and calling code should handle errors by
	// explicitly testing whether the error equals nil.

	// Ex.
	val, errAtoi := strconv.Atoi("42")

	// A nil error denotes success; a non-nil error denotes failure.
	if errAtoi != nil {
		fmt.Printf("couldn't convert number: %v\n", errAtoi)
		return
	}
	fmt.Println("Converted integer:", val)

	// See https://go.dev/tour/methods/19
	// for how to create own formatted error using same method as Stringer interface implementation

	// --- Readers ---

	// In the io package
	// Go std library contains many implementations of this interface,
	// including files, network connections, compressors, ciphers and others

	// The io.Reader interface has a Read method:
	// func (T) Read(b []byte) (n int, err error)

	// This method populates a given byte slice with data, then returns the num bytes populated (n) and an error value
	// Returns an io.EOF error when the stream ends

	// Ex. implemenation of the Reader interface - strings package
	reader := strings.NewReader("Hello, Reader")

	// set the max num bytes that can be read per iteration to 8 (equal to the set length of the slice)
	var bytes []byte = make([]byte, 8)

	// like StringBuilder in Java
	var readValueBdr strings.Builder

	var err error
	var n int

	// while reader not at EOF
	for err != io.EOF {
		n, err = reader.Read(bytes)
		readValueBdr.Write(bytes[:n])

		fmt.Printf("n = %v err = %v b = %v\n", n, err, bytes)
		fmt.Printf("b[:n] = %q\n", bytes[:n])
	}

	fmt.Println(readValueBdr.String())
}

func main() {
	// methodExamples()
	interfaceExamples()
}
