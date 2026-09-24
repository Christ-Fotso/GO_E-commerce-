package console

import "fmt"

func DisplayResult(operation string, result float64) {
	fmt.Printf("Le résultat pour l'opération %s est %.2f\n", operation, result)
}
