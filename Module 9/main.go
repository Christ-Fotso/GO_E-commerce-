package main

import (
	"fmt"

	"github.com/aminnairi/module8/library/algebra"
	"github.com/aminnairi/module8/library/console"
)

func main() {
	operation, numbers := console.ParseArguments()

	switch operation {
	case "add":
		console.DisplayResult("add", algebra.Add(numbers))

	case "sub":
		console.DisplayResult("sub", algebra.Subtract(numbers))

	case "mul":
		console.DisplayResult("mul", algebra.Multiply(numbers))

	default:
		fmt.Printf("Erreur : opération non reconnue (%s), attendu : add, mul, sub.", operation)
	}
}
