package main

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
)

func handleButton(mCache *measurementsCache) {
	sw1Code := 18
	fmt.Println("Button: Initialising button on line", sw1Code)
	sw1Line, err := gpiod.RequestLine("gpiochip0", sw1Code,
		gpiod.WithFallingEdge,
		gpiod.WithDebounce(time.Millisecond*200),
		gpiod.WithEventHandler(func(le gpiod.LineEvent) {
			fmt.Println("BUT: Received event:", le)
			mCache.clear()
		}))

	if err != nil {
		fmt.Println("BUT: Error requesting line:", err)
	}
	fmt.Println("BUT: Initialising... Done")
	defer sw1Line.Close()

	// block forever
	<-make(chan bool)
}
