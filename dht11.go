package main

import (
	"fmt"
	"time"
)

type measurement struct {
	temperature int8
	humidity    int8
	time        time.Time
	err         error
}

func handleDht(cache *measurementsCache, channel chan measurement) {
	for {
		fmt.Println("DHT: Sleeping for 1 second...")
		time.Sleep(time.Second * 1)
		fmt.Println("DHT: Sleeping for 1 second... Done")
		m := readFromDht()
		fmt.Printf("DHT: Temperature: %d, Humidity: %d, Time: %s\n", m.temperature, m.humidity, m.time)

		cache.addFront(m)
		channel <- m
	}
}

func readFromDht() measurement {
	fmt.Println("DHT: Reading from Dht...")

	time.Sleep(time.Second * 1)

	fmt.Println("DHT: Reading from Dht... Done successfully")
	return measurement{
		temperature: 21,
		humidity:    45,
		time:        time.Now(),
		err:         nil,
	}
}
