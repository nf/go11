package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go say(wg, "let's go!", 3)
	go say(wg, "ho!", 2)
	go say(wg, "hey!", 1)
	wg.Wait()
}

func say(wg *sync.WaitGroup, text string, secs int) {
	time.Sleep(time.Duration(secs) * time.Second)
	fmt.Println(text)
	wg.Done()
}
