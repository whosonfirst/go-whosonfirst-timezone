package timezones

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"testing"
)

func TestPointInPolygon(t *testing.T) {

	ctx := context.Background()

	tests := map[string][2]float64{
		"102047421": [2]float64{-122.384048, 37.616951},
	}

	schemes := []string{
		"timezones://",
		"timezones://github",
	}

	for _, uri := range schemes {

		l, err := NewTimezoneLookup(ctx, uri)

		if err != nil {
			t.Fatalf("Failed to create timezone lookup for %s, %v", uri, err)
		}

		for expected_id, coords := range tests {

			c, err := geo.NewCoordinate(coords[0], coords[1])

			if err != nil {
				t.Fatalf("Failed to create coordinate for %v, %v", coords, err)
			}

			rsp, err := l.PointInPolygon(ctx, c)

			if err != nil {
				t.Fatalf("Failed to perform point in polygon, %v", err)
			}

			results := rsp.Results()
			count := len(results)

			if count != 1 {
				t.Fatalf("Expected a single result, got %d", count)
			}

			first := results[0]

			if first.Id() != expected_id {
				t.Fatalf("Expected ID %s but got %s", expected_id, first.Id())
			}
		}
	}
}
