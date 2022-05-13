package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	MinTemperature int8
	MaxTemperature int8
	MinHumidity    int8
	MaxHumidity    int8
	Times          []string
	Temperatures   []string
	Humidities     []string
}

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

	data := generatePageData(mCache, bCache)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		fmt.Println("SER: Error parsing index.html:", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		fmt.Println("SER: Error executing template:", err)
		return
	}
}

func generatePageData(mCache *measurementsCache, bCache *boundsCache) PageData {
	measurements := mCache.getArray()
	times := []string{}
	temperatures := []string{}
	humidities := []string{}

	for _, m := range measurements {
		if m.err != nil {
			continue
		}

		times = append(times, m.time.Format("15:04:05"))
		temperatures = append(temperatures, fmt.Sprintf("%d", m.temperature))
		humidities = append(humidities, fmt.Sprintf("%d", m.humidity))
	}

	boundries := bCache.get()

	return PageData{
		Times:          times,
		Temperatures:   temperatures,
		Humidities:     humidities,
		MinTemperature: boundries.minTemperature,
		MaxTemperature: boundries.maxTemperature,
		MinHumidity:    boundries.minHumidity,
		MaxHumidity:    boundries.maxHumidity,
	}
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

	minTemperatureInt, err := parseToInt8(minTemperature)
	if err != nil {
		fmt.Println("SER: Error parsing minTemperature:", err)
		fmt.Fprint(w, "Error parsing minTemperature:", err)
		return
	}
	maxTemperatureInt, err := parseToInt8(maxTemperature)
	if err != nil {
		fmt.Println("SER: Error parsing maxTemperature:", err)
		fmt.Fprint(w, "Error parsing maxTemperature:", err)
		return
	}
	minHumidityInt, err := parseToInt8(minHumidity)
	if err != nil {
		fmt.Println("SER: Error parsing minHumidity:", err)
		fmt.Fprint(w, "Error parsing minHumidity:", err)
		return
	}
	maxHumidityInt, err := parseToInt8(maxHumidity)
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

func parseToInt8(s string) (int8, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}
