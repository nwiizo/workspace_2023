package main

import "fmt"

// fizzbuzz prints the numbers from 1 to 100. But for multiples of three it
// prints "Fizz" instead of the number and for the multiples of five it prints
// "Buzz". For numbers which are multiples of both three and five it prints
// "FizzBuzz".
func fizzbuzz() {
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}

func main() {
	fizzbuzz()
	fmt.Println("Hello, world!")
}
