# Test file example

Here provided the test `go` file example with mashed imports (including edge cases 
like dot imports, side-effect imports, import aliases).

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	. "math" // Dot import for direct access to exported identifiers
	_ "net/http/pprof" // Side-effect import for automatically registering pprof handlers

	pretty "github.com/kr/pretty" // External package with alias

	"github.com/wundersasha/sort-imports/utils" // Internal package
)

func main() {
	// Using a standard library import
	fmt.Println("Starting the program...")

	// Using a dot import
	fmt.Printf("The square root of 16 is: %v\n", Sqrt(16))

	// Using an external package to generate a UUID
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	fmt.Printf("Generated UUID: %s\n", id)

	// Using an external package with an alias
	fmt.Printf("Pretty printing with an alias: %s\n", pretty.Sprint(time.Now()))

	// Using an internal package function
	utils.DoSomething()

	// The side-effect import (_ "net/http/pprof") is not directly used but assumed to perform some initialization

	// Exit the program
	fmt.Println("Exiting the program...")
	os.Exit(0)
}
```