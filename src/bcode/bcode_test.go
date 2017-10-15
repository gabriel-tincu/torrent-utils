package bcode

import (
	"fmt"
	"runtime/debug"
	"testing"
)

var unmarshalTestData = []byte("d4:infod6:lengthi170917888e12:piece lengthi262144e4:name30:debian-8.8.0-arm64-netinst.isoe8:announce38:udp://tracker.publicbt.com:80/announce13:announce-listll38:udp://tracker.publicbt.com:80/announceel44:udp://tracker.openbittorrent.com:80/announceee7:comment33:Debian CD from cdimage.debian.orge")
var marshalTestData = Bencoded{
	"announce": []byte("udp://tracker.publicbt.com:80/announce"),
	"announce-list": []interface{}{
		[]interface{}{[]byte("udp://tracker.publicbt.com:80/announce")},
		[]interface{}{[]byte("udp://tracker.openbittorrent.com:80/announce")},
	},
	"comment": []byte("Debian CD from cdimage.debian.org"),
	"info": map[string]interface{}{
		"name":         []byte("debian-8.8.0-arm64-netinst.iso"),
		"length":       170917888,
		"piece length": 262144,
	},
}

func TestDecodeInt(t *testing.T) {
	res, _, err := decodeInt([]byte("i42eeeee"))
	if err != nil {
		t.Fatal(err)
	}
	if res != 42 {
		t.Fatal("result should be equal to 42")
	}
}

func TestDecodeString(t *testing.T) {
	//res, _, err := decodeByte([]byte)
}

func TestMarshal(t *testing.T) {
	data, err := Marshal(marshalTestData)
	if err != nil {
		t.Errorf("error marshaling: %s", err)
	}
	t.Log(fmt.Sprintf("got marshal result %+v", data))
}

func TestUnmarshal(t *testing.T) {
	_, err := Unmarshal([]byte(unmarshalTestData))
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(marshalTestData)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Unmarshal(unmarshalTestData)
	}
}
