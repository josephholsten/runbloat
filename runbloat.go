package main

// RunBloat takes a GPX file (from a run tracking app) and adds a little bit of random noise to the GPS locations. This makes it have a slightly larger total distance, while still maintaining the rough course and time.

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

var (
	filePath    = flag.String("f", "", "GPX file path")
	scaleFactor = flag.Float64("s", 1, "Fuzz scale factor")
)

func main() {
	var err error
	var gpxData *gpx.GPX

	flag.Parse()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Opening GPX file:", "foo", *filePath)
	if gpxData, err = gpx.ParseFile(*filePath); err != nil {
		fmt.Println("Could not parse GPX file: ", err)
		os.Exit(1)
	}

	var lastPoint gpx.GPXPoint
	var fuzzLat float64
	var fuzzLon float64
	// Analyize/manipulate your track data here...
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for i, _ := range segment.Points {
				point := &(segment.Points[i])
				if lastPoint.Latitude != 0 {
					fuzzLat = *scaleFactor * r.Float64() * (lastPoint.Latitude - point.Latitude)
					fuzzLon = *scaleFactor * r.Float64() * (lastPoint.Longitude - point.Longitude)
				}
				// fmt.Printf("%v, %v, %v, %v\n", point.Latitude, fuzzLat, point.Longitude, fuzzLon)

				lastPoint = *point

				point.Latitude += fuzzLat
				point.Longitude += fuzzLon
				// fmt.Printf("%v, %v, %v, %v\n", point.Latitude, fuzzLat, point.Longitude, fuzzLon)
			}
		}
	}

	// fmt.Print(gpx.GetGpxElementInfo("", gpxData))

	// (Check the API for GPX manipulation and analyzing utility methods)

	// When ready, you can write the resulting GPX file:
	xmlBytes, err := gpxData.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
	if err != nil {
		fmt.Println("Could not convert GPX to XML: ", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", xmlBytes)
}
