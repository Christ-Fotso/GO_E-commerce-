package algebra

import "testing"

func TestAddTwoNumbers(t *testing.T) {
	result := Add([]float64{1, 2})
	var expectation float64 = 3

	if result != expectation {
		t.Errorf("Expected %f to equal %f\n", result, expectation)
	}
}
