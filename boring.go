package main

import (
	"fmt"
	"math/rand"
)

func multiplex(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				// do nothing
			case <-quit:
			  // do cleanup stuff here
			  quit <- "See you later!"
				return
			}
		}
	}()

	return c
}

func main() {
	quit := make(chan string)
	c := boring("Han", quit)

	for i := rand.Intn(30); i >= 0; i-- {
		fmt.Println(<-c)
	}
	quit <- "Bye!"
  fmt.Printf("Han says: %q\n", <-quit)
}
