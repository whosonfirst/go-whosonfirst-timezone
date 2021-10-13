package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-timezones"
	"log"
	"os"
	"time"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "timezones://", "...")

	latitude := flag.Float64("latitude", 0.0, "")
	longitude := flag.Float64("longitude", 0.0, "")

	flag.Parse()

	ctx := context.Background()

	t1 := time.Now()
	tz_lookup, err := timezones.NewTimezoneLookup(ctx, *lookup_uri)

	log.Printf("Time to compile timezone lookup, %v", time.Since(t1))

	if err != nil {
		log.Fatalf("Failed to create new timezones lookup, %v", err)
	}

	// Something something something... go-whosonfirst-spatial-pip.PointInPolygonRequest

	c, err := geo.NewCoordinate(*longitude, *latitude)

	if err != nil {
		log.Fatalf("Failed to create coordinate, %v", err)
	}

	// Something something something... go-whosonfirst-spatial-pip filters...
	// f, err := NewSPRFilterFromPointInPolygonRequest(req)

	rsp, err := tz_lookup.PointInPolygon(ctx, c)

	if err != nil {
		log.Fatalf("Failed to perform point in polygon request, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(rsp)

	if err != nil {
		log.Fatalf("Failed to encode results, %v", err)
	}
}
