package exampleNoSecret

import "fmt"

// Go source code that will not match the detection rule

func anotherFunction() {
	fmt.Printf("Not a database here")
}

func main() {
	fmt.Printf("Definitely no secrets in this code")
}
