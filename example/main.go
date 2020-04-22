package main

import (
	"fmt"
	"log"

	"github.com/abibby/gpmdp"
)

func main() {
	g, err := gpmdp.Connect()
	if err != nil {
		log.Fatal(err)
	}
	track := &gpmdp.Track{}
	playing := false
	for {
		select {
		case err := <-g.Error:
			log.Fatal(err)
		case ev := <-g.Event:
			if t, ok := ev.Track(); ok {
				track = t
			}
			if p, ok := ev.Playing(); ok {
				playing = p
			}
			fmt.Printf("%5v %s - %s\n", playing, track.Title, track.Artist)
		}
	}
}
