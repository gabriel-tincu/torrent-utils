package magnet

import (
	"io/ioutil"
	"testing"
)


var magnetAddr = "magnet:?xt=urn:btih:7fbc58e324b539bdda58c15bda3acd26b0d5fbd1&dn=Luis+Fonsi+-+Despacito+%28feat.+Daddy+Yankee%29&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"

func TestMagnet(t *testing.T) {
	bits, err := ioutil.ReadFile("/home/gabi/.config/transmission/torrents/BrandonSanderson.works.62f7b36331b00536.torrent")
	_ = bits
	if err != nil {
		t.Fatal(err)
	}
	res, err := ParseMagnet([]byte(magnetAddr))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res.Trackers {
		t.Log(v)
	}
	t.Log(res.DisplayName)
}
