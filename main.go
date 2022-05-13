package main

func main() {
	mCache := measurementsCacheInit()
	boundsCache := boundsCacheInit()
	channel := make(chan measurement)

	go handleDht(mCache, channel)
	go handleLeds(boundsCache, channel)
	go handleButton(mCache)

	handleHttpServer(mCache, boundsCache)
}
