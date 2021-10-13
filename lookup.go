package timezones

import (
	"context"
	"fmt"
	"github.com/aaronland/go-jsonl/walk"
	"github.com/paulmach/orb"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-timezones/data"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"sync"
)

var lookup_table database.SpatialDatabase
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type TimezoneLookupFunc func(context.Context)

type TimezoneLookup struct {
}

func NewTimezoneLookup(ctx context.Context, uri string) (*TimezoneLookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var source string

	switch u.Host {
	case "timezone":
		source = u.Path
	default:
		source = u.Host
	}

	switch source {
	case "iterator":

		return nil, fmt.Errorf("Not implemented")

	case "github":

		rsp, err := http.Get(DATA_GITHUB)

		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve data from %s, %w", DATA_GITHUB, err)
		}

		lookup_func := NewTimezoneLookupFuncWithReader(ctx, rsp.Body)
		return NewTimezoneLookupWithLookupFunc(ctx, lookup_func)

	default:

		fs := data.FS
		fh, err := fs.Open(DATA_JSON)

		if err != nil {
			return nil, fmt.Errorf("Failed to load local precompiled data, %w", err)
		}

		lookup_func := NewTimezoneLookupFuncWithReader(ctx, fh)
		return NewTimezoneLookupWithLookupFunc(ctx, lookup_func)
	}
}

func NewTimezoneLookupFuncWithReader(ctx context.Context, r io.ReadCloser) TimezoneLookupFunc {

	lookup_func := func(ctx context.Context) {

		defer r.Close()

		db, err := database.NewSpatialDatabase(ctx, "sqlite://?dsn=:memory:")

		if err != nil {
			lookup_init_err = fmt.Errorf("Failed to initialize database, %w", err)
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var walk_err error

		record_ch := make(chan *walk.WalkRecord)
		error_ch := make(chan *walk.WalkError)
		done_ch := make(chan bool)

		go func() {

			for {
				select {
				case <-ctx.Done():
					done_ch <- true
					return
				case err := <-error_ch:
					walk_err = err
					done_ch <- true
				case r := <-record_ch:

					err := db.IndexFeature(ctx, r.Body)

					if err != nil {
						error_ch <- &walk.WalkError{
							Path:       r.Path,
							LineNumber: r.LineNumber,
							Err:        fmt.Errorf("Failed to index feature, %w", err),
						}
					}
				}
			}
		}()

		walk_opts := &walk.WalkOptions{
			IsBzip:        true,
			RecordChannel: record_ch,
			ErrorChannel:  error_ch,
			Workers:       10,
		}

		walk.WalkReader(ctx, walk_opts, r)

		<-done_ch

		if walk_err != nil && !walk.IsEOFError(walk_err) {
			lookup_init_err = walk_err
			return
		}

		lookup_table = db
	}

	return lookup_func
}

func NewTimezoneLookupWithLookupFunc(ctx context.Context, lookup_func TimezoneLookupFunc) (*TimezoneLookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := TimezoneLookup{}
	return &l, nil
}

func (l *TimezoneLookup) PointInPolygon(ctx context.Context, coord *orb.Point) (spr.StandardPlacesResults, error) {
	return lookup_table.PointInPolygon(ctx, coord)
}
