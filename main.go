package main

import (
	"fmt"
	"monkey/repl"
)

func main() {
	fmt.Printf("Starting REPL...\n----------------\n\n")
	// repl.Start()
	repl.File("test.monk")
}
