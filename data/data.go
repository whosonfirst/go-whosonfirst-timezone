package data

import (
	"embed"
)

/*
$> cd /usr/local/whosonfirst-go-exportify
$> ./bin/as-jsonl -iterator-uri 'repo://?include=properties.wof:placetype=timezone&include=properties.mz:is_current=1' /usr/local/data/whosonfirst-data-admin-xy/ > timezones.jsonl
2021/10/13 09:50:31 time to index paths (1) 2.512386626s
*/

//go:embed *.jsonl.bz2
var FS embed.FS
