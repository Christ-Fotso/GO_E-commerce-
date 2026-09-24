package algebra

import "testing"

func TestMultiplyTwoNumbers(t *testing.T) {
	a := 12.5
	b := 7.0
	result := Multiply([]float64{a, b})
	expectation := 87.5

	if result != expectation {
		t.Errorf("Multiply(%f, %f) should be %f, received %f\n", a, b, expectation, result)
	}
}
