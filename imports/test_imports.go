package imports

import (
	"fmt"
	
	"github.com/google/uuid"
	
	"github.com/wundersasha/sort-imports/utils"
)

func someTestFunction() {
	fmt.Println("Hello, World!")
	fmt.Println(uuid.New())
	utils.DoSomething()
}
