package algebra

func Subtract(numbers []float64) (subtraction float64) {
	for _, number := range numbers {
		subtraction -= number
	}

	return
}
