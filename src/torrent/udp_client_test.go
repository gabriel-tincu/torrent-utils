package torrent

import (
	"testing"
)

func TestConnect(t *testing.T) {
	data, connectId := getConnectBytes()
	resp, err := sendUDP("tracker.leechers-paradise.org:6969", data)
	if err != nil {
		t.Fatal(err)
	}
	transaction, err := parseConnectResponse(resp, connectId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(transaction)
}
