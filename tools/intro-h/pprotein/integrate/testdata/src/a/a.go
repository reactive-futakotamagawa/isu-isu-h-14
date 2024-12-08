package a

import ( // want "change import"
	"fmt"
	"math"
)

func f() {
	// The pattern can be written in regular expression.
	var gopher int
	print(gopher)
}

func main() { // want "main"
	f()
	fmt.Printf("Hello, World! %d", math.MaxInt64)
}
