package main

import (
	"fmt"
	"math/rand"
	"time"
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

func boring(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c
}

func main() {
	c := boring("Han")
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-time.After(500 * time.Millisecond):
			fmt.Println("You're too slow.")
			return
		}
	}
}
