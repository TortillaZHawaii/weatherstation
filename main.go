package main

func main() {
	mCache := measurementsCacheInit()
	boundsCache := boundsCacheInit()

	go handleDht(mCache)
	// go handleButton(&mCache)

	handleHttpServer(mCache, boundsCache)
}
