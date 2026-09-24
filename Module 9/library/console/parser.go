package console

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func ParseArguments() (string, []float64) {
	flag.Parse()

	args := flag.Args()
	
	operation := args[0]
	texts := args[1:]
	numbers := make([]float64, len(texts))

	for index, text := range texts {
		number, err := strconv.ParseFloat(text, 0)

		if err != nil {
			fmt.Printf("Erreur : L'argument #%d n'est pas un nombre, %v reçu.", index, text)
			os.Exit(1)
		}

		numbers = append(numbers, number)
	}

	return operation, numbers
}
