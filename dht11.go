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
	temperatureI32, err := strconv.Atoi(string(tempFile))

	if err != nil {
		return measurement{err: err}
	}

	humFile, err := os.ReadFile("/sys/bus/iio/devices/iio:device0/in_humidityrelative_input")
	if err != nil {
		return measurement{err: err}
	}

	fmt.Println("DHT: Read humidity", string(humFile))
	humI32, err := strconv.Atoi(string(humFile))

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

func retrieveMeasurementFromBits(bits []byte) measurement {
	s := string(bits)
	fmt.Println("DHT: Parsing bits:", s)
	temperature, err := strconv.ParseInt(s[:8], 2, 8)
	if err != nil {
		return measurement{err: err}
	}
	humidity, err := strconv.ParseInt(s[16:24], 2, 8)
	if err != nil {
		return measurement{err: err}
	}

	checksum, err := strconv.ParseInt(s[32:40], 2, 8)
	if err != nil {
		return measurement{err: err}
	}

	fmt.Println("DHT: Checking checksum...")
	checksumCalculated := 0
	for _, a := range s {
		if a == '1' {
			checksumCalculated++
		}
	}

	if checksumCalculated != int(checksum) {
		return measurement{err: fmt.Errorf("checksum failed")}
	}

	return measurement{
		temperature: int8(temperature),
		humidity:    int8(humidity),
		time:        time.Now(),
		err:         nil,
	}
}

func failAfterTime(blockChan chan bool, duration time.Duration) {
	time.Sleep(duration)
	blockChan <- false
}
