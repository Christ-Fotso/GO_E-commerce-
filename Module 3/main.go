package main

import (
	"errors"
	"fmt"
)

func add(a, b float64) float64 {
	return a + b
}

func subtract(a, b float64) float64 {
	return a - b
}

func multiply(a, b float64) float64 {
	return a * b
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Impossible de diviser par zéro")
	}

	return a / b, nil
}

func surround(prefix, suffix string) func(string) string {
	return func(text string) string {
		return fmt.Sprintf("%s%s%s", prefix, text, suffix)
	}
}

func main() {
	fmt.Println("Résultat de l'addition", add(1, 2))
	fmt.Println("Résultat de la soustraction", subtract(1, 2))
	fmt.Println("Résultat de la multiplication", multiply(1, 2))

	result, err := divide(1, 2)

	if err != nil {
		fmt.Printf("Erreur lors de la division : %s\n", err.Error())
	} else {
		fmt.Println("Résultat de la division", result)
	}

	surroundWithDoubleBrackets := surround("[[", "]]")

	fmt.Println(surroundWithDoubleBrackets("Go"))
	fmt.Println(surroundWithDoubleBrackets("ESGI"))
	fmt.Println(surroundWithDoubleBrackets("2026"))
}
