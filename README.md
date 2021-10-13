# go-whosonfirst-timezone


## Important

This work is still experimental. It will probably change. It may be scrapped entirely.

## Example

```
> go run -mod vendor cmd/lookup/main.go -latitude 37.616951 -longitude -122.384048 | jq
2021/10/13 10:25:59 Time to compile timezone lookup, 11.305766815s
{
  "places": [
    {
      "wof:id": "102047421",
      "wof:parent_id": "85633793",
      "wof:name": "America/Los_Angeles",
      "wof:country": "",
      "wof:placetype": "timezone",
      "mz:latitude": 38.27008,
      "mz:longitude": -118.219968,
      "mz:min_latitude": 32.534622,
      "mz:min_longitude": 38.27008,
      "mz:max_latitude": -124.733253,
      "mz:max_longitude": -114.039345,
      "mz:is_current": -1,
      "mz:is_deprecated": 0,
      "mz:is_ceased": -1,
      "mz:is_superseded": 0,
      "mz:is_superseding": 0,
      "edtf:inception": "",
      "edtf:cessation": "",
      "wof:supersedes": [],
      "wof:superseded_by": [],
      "wof:belongsto": [
        102191575,
        85633793
      ],
      "wof:path": "102/047/421/102047421.geojson",
      "wof:repo": "whosonfirst-data-admin-xy",
      "wof:lastmodified": 1566655657
    }
  ]
}
```