package algebra

func Multiply(numbers []float64) (float64) {
	var multiplication float64 = 1

	for _, number := range numbers {
		multiplication *= number
	}

	return multiplication
}
