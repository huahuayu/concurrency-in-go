package main

import (
	"fmt"
)

func main(){
	sayHello := func() {
		fmt.Println("hello from goroutin")
	}

	go sayHello()

	fmt.Println("hello from main")
}
