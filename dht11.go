package main

import (
	"fmt"
	"os"
	"strconv"
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
		sleepSecondsTime := 5
		fmt.Println("DHT: Sleeping for", sleepSecondsTime, "seconds...")
		time.Sleep(time.Duration(time.Second * 5))
		fmt.Println("DHT: Sleeping for", sleepSecondsTime, "seconds... Done")
		m := readFromDht()
		fmt.Printf("DHT: Temperature: %d, Humidity: %d, Time: %s\n", m.temperature, m.humidity, m.time)

		cache.addFront(m)
		channel <- m
	}
}

func readFromDht() measurement {
	fmt.Println("DHT: Reading from Dht...")

	tempFile, err := os.ReadFile("/sys/bus/iio/devices/iio:device0/in_temp_input")
	if err != nil {
		return measurement{err: err}
	}

	fmt.Println("DHT: Read temp", string(tempFile))
	tempStr := string(tempFile)
	tempStr = tempStr[:len(tempStr)-1]
	temperatureI32, err := strconv.Atoi(tempStr)

	if err != nil {
		return measurement{err: err}
	}

	humFile, err := os.ReadFile("/sys/bus/iio/devices/iio:device0/in_humidityrelative_input")
	if err != nil {
		return measurement{err: err}
	}

	fmt.Println("DHT: Read humidity", string(humFile))
	humStr := string(humFile)
	humStr = humStr[:len(humStr)-1]
	humI32, err := strconv.Atoi(humStr)

	if err != nil {
		return measurement{err: err}
	}

	return measurement{
		temperature: int8(temperatureI32 / 1000),
		humidity:    int8(humI32 / 1000),
		time:        time.Now(),
		err:         nil,
	}
}
