package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/warthog618/gpiod"
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

	dht, err := gpiod.RequestLine("gpiochip0", 4,
		gpiod.AsOutput(1))

	if err != nil {
		fmt.Println("DHT: Error requesting line:", err)
		return measurement{err: err}
	}

	blockChan := make(chan bool)
	go failAfterTime(blockChan, time.Second*2)

	// start signal
	dht.SetValue(0)
	time.Sleep(time.Millisecond * 18)
	dht.SetValue(1)
	dht.Close()

	i := 0
	start := time.Duration(0)

	// wait for response

	bits := make([]byte, 50)

	dht, err = gpiod.RequestLine("gpiochip0", 4,
		gpiod.WithBothEdges,
		gpiod.WithEventHandler(func(le gpiod.LineEvent) {
			if le.Type == gpiod.LineEventRisingEdge {
				start = le.Timestamp
			} else if le.Type == gpiod.LineEventFallingEdge {
				end := le.Timestamp
				diff := end - start
				i++

				if diff > time.Microsecond*80 { // 80 mikros start streaming
					i = -1
				} else if diff > time.Microsecond*64 { // 70 mikros high
					bits[i] = '1'
				} else if diff < time.Microsecond*30 { // 24 mikros low
					bits[i] = '0'
				} else if diff < time.Microsecond*15 { // error
				}

				fmt.Println("DHT: i:", i, "Diff:", diff)

				allSignalsReceived := i == 39
				if allSignalsReceived {
					blockChan <- true
				}
			}
		}))
	if err != nil {
		fmt.Println("DHT: Error requesting line:", err)
		return measurement{err: err}
	}
	defer dht.Close()

	// block
	isSuccessful := <-blockChan

	if isSuccessful {
		fmt.Println("DHT: Reading from Dht... Done")
		return retrieveMeasurementFromBits(bits)
	} else {
		fmt.Println("DHT: Reading from Dht... Failed")
		return measurement{err: fmt.Errorf("failed to read from DHT, wrong timing")}
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
