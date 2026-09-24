package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width float64
	Height float64
}

func (rectangle Rectangle) Area() float64 {
	return rectangle.Height * rectangle.Width
}

func (rectangle Rectangle) Perimeter() float64 {
	return rectangle.Area() * 2
}

type Circle struct {
	Radius float64
}

func (circle Circle) Perimeter() float64 {
	return 2 * math.Pi * circle.Radius
}

func (circle Circle) Area() float64 {
	return math.Pi * math.Pow(2, circle.Radius)
}

func details(shape Shape) string {
	return fmt.Sprintf("Rayon : %.2f — Périmètre : %.2f\n", shape.Area(), shape.Perimeter())
}

func main() {
	rectangle := Rectangle{
		Width: 5.2,
		Height: 10.4,
	}	

	circle := Circle{
		Radius: 7.5,
	}

	fmt.Println(details(rectangle))
	fmt.Println(details(circle))
}
