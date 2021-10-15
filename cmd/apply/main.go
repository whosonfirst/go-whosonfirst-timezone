// apply is a command-line tool to "apply" timezone information to a record. As written it is doing a point-in-polygon
// query which is not suitable for large (area) features that span multiple timezones. For example "Asia". TBD...
package main

/*

go run -mod vendor cmd/apply/main.go -iterator-uri 'repo://?exclude=properties.wof:placetype=timezone' -writer-uri fs:///usr/local/data/whosonfirst-data-admin-xy/data /usr/local/data/whosonfirst-data-admin-xy/

*/

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-timezones"
	"github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
	"io"
	"log"
	"strconv"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "timezones://", "A valid whosonfirst/go-whosonfirst-timezones URI.")
	iterator_uri := flag.String("iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2 URI.")
	writer_uri := flag.String("writer-uri", "stdout://", "A valid whosonfirst/go-writer URI.")

	flag.Parse()

	iterator_sources := flag.Args()

	ctx := context.Background()

	tz_lookup, err := timezones.NewTimezoneLookup(ctx, *lookup_uri)

	if err != nil {
		log.Fatalf("Failed to create new timezones lookup, %v", err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create new writer, %v", err)
	}

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse URI for %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		// START OF something something something intersects
		// https://pkg.go.dev/github.com/go-spatial/geom@v0.0.0-20210804045141-e17f18e24cae/planar/intersect
		
		pt, _, err := properties.Centroid(body)

		if err != nil {
			return fmt.Errorf("Failed to derive centroid for %s, %w", path, err)
		}

		rsp, err := tz_lookup.PointInPolygon(ctx, pt)

		if err != nil {
			return fmt.Errorf("Failed to perform point in polygon request for %s, %v", path, err)
		}

		// END OF something something something intersects

		results := rsp.Results()
		count := len(results)

		timezone_ids := make([]int64, count)

		for idx, r := range results {

			id, err := strconv.ParseInt(r.Id(), 10, 64)

			if err != nil {
				return fmt.Errorf("Failed to derive wof:id for %s (%s), %w", path, r.Id(), err)
			}

			timezone_ids[idx] = id
		}

		to_update := map[string]interface{}{
			"properties.wof:timezones": timezone_ids,
		}

		changed, body, err := export.AssignPropertiesIfChanged(ctx, body, to_update)

		if err != nil {
			return fmt.Errorf("Failed to assign timezones to %s, %w", path, err)
		}

		if !changed {
			return nil
		}

		err = wof_writer.WriteFeatureBytes(ctx, wr, body)

		if err != nil {
			return fmt.Errorf("Failed to write %s, %w", path, err)
		}

		log.Printf("Updated %s\n", path)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create new iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		log.Fatalf("Failed to iterate sources, %v", err)
	}
}
