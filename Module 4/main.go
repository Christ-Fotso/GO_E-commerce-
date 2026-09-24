package main

import (
	"fmt"
)

func studentAverage(grades []float64) (average float64) {
	for _, grade := range grades {
		average += grade
	}	

	average /= float64(len(grades))

	return
}

func classAverage(class map[string][]float64) (average float64) {
	for _, grades := range class {
		average += studentAverage(grades)
	}

	average /= float64(len(class))

	return
}

func main() {
	students := map[string][]float64{
		"Alice": { 15.5, 12, 18 },
		"Bob": { 8, 14 },
		"Charlie": { 10, 11.5, 9.5, 16 },
	}

	fmt.Printf("Moyenne de la classe : %.2f\n", classAverage(students))
}
