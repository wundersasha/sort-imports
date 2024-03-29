package imports

import (
	"fmt"
	"log"
	. "math"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/google/uuid"
	pretty "github.com/kr/pretty"

	"github.com/wundersasha/sort-imports/utils"
)

// Dot import for direct access to exported identifiers
// Side-effect import for automatically registering pprof handlers

// External package with alias

// Internal package

func someFunc() {
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
