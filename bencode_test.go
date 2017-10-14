package bencode

import (
	"fmt"
	"net/url"
	"testing"
)

var magnet = "magnet:?xt=urn:btih:be257fdc49adee6bbc6dbc18445f8313f7916c3e&dn=IT+2017+Movies+HD+TS+XviD+Clean+English+Audio+New+%2BSample+%E2%98%BBrDX%E2%98%BB&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"

func TestEscape(t *testing.T) {

	parsed, _ := url.ParseQuery(magnet[8:])
	for key, val := range parsed {
		fmt.Printf("k: %s\tv: %+v\n", key, val)
	}
}
