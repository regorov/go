package main

import (
	"fmt"
	"github.com/kierdavis/go/chem"
	"os"
)

func main() {
	c := chem.ParseCompound(os.Args[1])

	fmt.Printf("Name: %s\n", c.Name())
	fmt.Printf("Formala: %s\n", c.String())
	fmt.Printf("Relative molecular mass: %.1f\n", c.RoundMass())
}
