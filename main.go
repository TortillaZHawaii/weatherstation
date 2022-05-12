package main

import (
	"fmt"
	"time"
)

func httpServer() {
	fmt.Println("Server!")
}

func handleLeds() {
	fmt.Println("Leds")
}

func handleButton() {
	fmt.Println("Button")
}

func readFromDht() {
	fmt.Println("Dht")
}

func main() {
	go httpServer()
	go readFromDht()
	go handleButton()
	go handleLeds()
	time.Sleep(time.Second * 2)
	fmt.Println("End")
}
