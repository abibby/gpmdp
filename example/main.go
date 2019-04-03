package main

import (
	"fmt"
	"log"

	"github.com/zwzn/gpmdp"
)

func main() {
	g, err := gpmdp.Connect()
	if err != nil {
		log.Fatal(err)
	}
	playing := false
	track := &gpmdp.Track{}
	for {
		select {
		case track = <-g.Track():
		case playing = <-g.Playing():
		}
		fmt.Printf("%5v %s - %s\n", playing, track.Title, track.Artist)
	}
}
