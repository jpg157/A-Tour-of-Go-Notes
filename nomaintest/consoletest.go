package main

import "fmt"

// Outside a function, every statement begins with a keyword
// (var, func, and so on), and so the := construct is not available.
var _, _ = fmt.Println("hello world")

func main() {}
