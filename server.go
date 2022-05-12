package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func handleHttpServer(mCache *measurementsCache, bCache *boundsCache) {
	fmt.Println("SER: Httpserver starting...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGet(w, r, mCache, bCache)
		}
		if r.Method == "POST" {
			handlePost(w, r, mCache, bCache)
		}
	})
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handleGet(w http.ResponseWriter, r *http.Request,
	mCache *measurementsCache, bCache *boundsCache) {
	http.ServeFile(w, r, "index.html")
}

func handlePost(w http.ResponseWriter, r *http.Request,
	mCache *measurementsCache, bCache *boundsCache) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	minTemperature := r.FormValue("mintemperature")
	maxTemperature := r.FormValue("maxtemperature")
	minHumidity := r.FormValue("minhumidity")
	maxHumidity := r.FormValue("maxhumidity")

	fmt.Println("SER: Received new bounds:",
		minTemperature, "-", maxTemperature, "C,",
		minHumidity, "-", maxHumidity, "%")

	minTemperatureInt, err := parseToInt16(minTemperature)
	if err != nil {
		fmt.Println("SER: Error parsing minTemperature:", err)
		fmt.Fprint(w, "Error parsing minTemperature:", err)
		return
	}
	maxTemperatureInt, err := parseToInt16(maxTemperature)
	if err != nil {
		fmt.Println("SER: Error parsing maxTemperature:", err)
		fmt.Fprint(w, "Error parsing maxTemperature:", err)
		return
	}
	minHumidityInt, err := parseToInt16(minHumidity)
	if err != nil {
		fmt.Println("SER: Error parsing minHumidity:", err)
		fmt.Fprint(w, "Error parsing minHumidity:", err)
		return
	}
	maxHumidityInt, err := parseToInt16(maxHumidity)
	if err != nil {
		fmt.Println("SER: Error parsing maxHumidity:", err)
		fmt.Fprint(w, "Error parsing maxHumidity:", err)
		return
	}

	bCache.set(bounds{
		minTemperature: minTemperatureInt,
		maxTemperature: maxTemperatureInt,
		minHumidity:    minHumidityInt,
		maxHumidity:    maxHumidityInt,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func parseToInt16(s string) (int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}
