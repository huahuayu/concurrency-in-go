package main

import (
	"fmt"
)

func main() {
	fibonacci := func() func() int {
		back1, back2 := -1, 1
		return func() int {
			back1, back2 = back2, (back1 + back2)
			return back2
		}
	}

	f := fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
