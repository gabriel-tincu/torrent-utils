package bcode

import (
	"io/ioutil"
	"runtime/debug"
	"testing"
)

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

func TestBDecode(t *testing.T) {
	bits, err := ioutil.ReadFile("data/test.torrent")
	if err != nil {
		t.Fatal(err)
	}
	_, err = BDecode(bits)
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}
}

func BenchmarkBDecode(b *testing.B) {
	bits, err := ioutil.ReadFile("data/test.torrent")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = BDecode(bits)
	}
}
