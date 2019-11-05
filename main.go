package main

import (
	"fmt"
	"math/rand"
	"stringutil"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(stringutil.Reverse("Hello, kimchi!"))
	fmt.Println("Hello, kimchi!")
	fmt.Println("My favorite number is: ", rand.Intn(10))
	sumOfSeven := 7
	refSeven := &sumOfSeven
	*refSeven = *refSeven + 1
	fmt.Printf("Sum of 3 and 4 is %d \n", sumOfSeven)
	fmt.Println("Rotated a, b, c looks like ")
	fmt.Println(rotate("a", "b", "c"))
	var a, b rune = 1, 2

	fmt.Println(a, b)
}

func sum(x, y int) int {
	return x + y
}

func rotate(x, y, z string) (a, b, c string) {
	a = y
	b = z
	c = x
	return
}
