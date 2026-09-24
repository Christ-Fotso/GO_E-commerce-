package main

import "fmt"

func main() {
	const admissionThreshold float64 = 10

	grade := 14.5
	integerGrade := int(grade)

	if integerGrade >= 16 {
		fmt.Println("Excellent")
	} else if integerGrade >= 14 {
		fmt.Println("Bien")
	} else if integerGrade >= 10 {
		fmt.Println("Passable")
	} else {
		fmt.Println("Insuffisant")
	}

	for number := 1; number <= 10; number++ {
		fmt.Println(number)
	}

	words := []string{"Go", "est", "génial"}

	for _, word := range words {
		fmt.Println(word)
	}
}
