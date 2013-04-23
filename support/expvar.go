package main

import (
	"expvar"
	"log"
	"net/http"
	"time"
)

func main() {
	count := expvar.NewInt("count")
	go func() {
		for {
			count.Add(1)
			time.Sleep(time.Second)
		}
	}()
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
