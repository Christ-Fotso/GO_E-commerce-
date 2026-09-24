package main

import (
	"errors"
	"fmt"
)

type DivisionError struct {
	Numerator float64
}

func (err DivisionError) Error() string {
	return fmt.Sprintf("cannot divide %f with zero\n", err.Numerator)
}

type NegativeNumeratorError struct {
	Value float64
}

func (err NegativeNumeratorError) Error() string {
	return fmt.Sprintf("numerator cannot be negative, %f provided", err.Value)
}

type NegativeDenominatorError struct {
	Value float64
}

func (err NegativeDenominatorError) Error() string {
	return fmt.Sprintf("denominator cannot be negative, %f provided", err.Value)
}

func divide(numerator, denominator float64) (result float64, err error) {
	if numerator < 0 {
		return 0, NegativeNumeratorError{
			Value: numerator,
		}
	}

	if denominator < 0 {
		return 0, NegativeDenominatorError{
			Value: denominator,
		}
	}

	if denominator == 0 {
		return 0, DivisionError{
			Numerator: numerator,
		}
	}

	return numerator / denominator, nil
}

func main() {
	var numerator float64 = 10
	var denominator float64 = 5

	result, err := divide(numerator, denominator)
	if err != nil {
		var negativeDenominatorError NegativeDenominatorError

		if errors.As(err, &negativeDenominatorError) {
			fmt.Printf("Le dénominateur est négatif (%f)\n", negativeDenominatorError.Value)
		} else {
			err.Value
			fmt.Println(err.Error())
		}
	} else {
		fmt.Printf("La division de %.2f par %.2f fait %.2f\n", numerator, denominator, result)
	}
}
