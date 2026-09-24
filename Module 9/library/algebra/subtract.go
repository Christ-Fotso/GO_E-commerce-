package algebra

func Subtract(numbers []float64) (subtraction float64) {
	subtraction = numbers[0]

	for _, number := range numbers[1:] {
		subtraction -= number
	}

	return
}
