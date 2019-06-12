package main

import (
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup

	greeters := []string{"alice","bob","eva","david","frank"}

	for _, greeter := range greeters{
		wg.Add(1)
		go greeting(&wg, greeter)
	}

	wg.Wait()
}

func greeting(wg *sync.WaitGroup, name string){
	defer wg.Done()
	fmt.Printf("greeting from %s!\n", name)
}
