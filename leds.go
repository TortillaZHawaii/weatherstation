package main

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
)

type led struct {
	line *gpiod.Line
	name string
}

func ledInit(line int, name string) *led {
	fmt.Println("LED: Initialising led on line", line, "with name", name)

	red, err := gpiod.RequestLine("gpiochip0", line, gpiod.AsOutput(0))
	if err != nil {
		fmt.Println("LED: Error requesting", name, "line:", err)
		panic(err)
	}

	fmt.Println("LED: Initialising", name, "... Done")
	return &led{
		line: red,
		name: name,
	}
}

func (led *led) close() {
	fmt.Println("LED: Closing led", led.name)
	led.line.Reconfigure(gpiod.AsInput)
	led.line.Close()
}

func (led *led) blink() {
	fmt.Println("LED: Blinking", led.name)
	led.line.SetValue(1)
	time.Sleep(time.Millisecond * 500)
	led.line.SetValue(0)
}

func handleLeds(bCache *boundsCache, channel chan measurement) {
	fmt.Println("LED: Starting...")
	fmt.Println("LED: Initialising...")

	red := ledInit(27, "red")
	defer red.close()
	yellow := ledInit(23, "yellow")
	defer yellow.close()
	blue := ledInit(24, "blue")
	defer blue.close()

	for m := range channel {
		if m.err != nil {
			fmt.Println("LED: Received with error:", m.err)
			go red.blink()
			continue
		}

		boundries := bCache.get()
		fmt.Printf("LED: Temperature: %d, Humidity: %d, Time: %s\n", m.temperature, m.humidity, m.time)
		if m.temperature < boundries.minTemperature || m.temperature > boundries.maxTemperature {
			fmt.Println("LED: Temperature is outside of bounds")
			go yellow.blink()
		}
		if m.humidity < boundries.minHumidity || m.humidity > boundries.maxHumidity {
			fmt.Println("LED: Humidity is outside of bounds")
			go blue.blink()
		}
	}
}
