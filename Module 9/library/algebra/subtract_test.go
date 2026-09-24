package algebra

import (
	"math"
	"testing"
)

func TestSubtractingTwoNumbers(t *testing.T) {
	numbers := []float64{12.7, 6.4, 3.2}
	expectation := 3.1
	result := Subtract(numbers)
	tolerance := 1e-9

	if math.Abs(result-expectation) > tolerance {
		t.Errorf("Subtract() should return %f for %v but %f received\n", expectation, numbers, result)
	}
}
