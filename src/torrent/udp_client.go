package torrent

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	magicNumber   = 4497486125440
	connectAction = 0
	udpNet        = "udp"
	buffSize      = 10000
	readDeadline  = time.Second * 2
)

func getConnectBytes() (response []byte, connectionId uint32) {
	response = make([]byte, 16)
	binary.BigEndian.PutUint64(response, magicNumber)
	binary.BigEndian.PutUint32(response[8:12], connectAction)
	connectionId = rand.Uint32()
	binary.BigEndian.PutUint32(response[12:], connectionId)
	return
}

func parseConnectResponse(data []byte,transactionId uint32) (connectionId []byte, err error) {
	if len(data) != 16 {
		err = fmt.Errorf("response size malformed :%d", len(data))
		return
	}
	action := binary.BigEndian.Uint32(data[:4])
	if  action != connectAction {
		err = fmt.Errorf("action id should be connect : %d", action)
		return
	}
	transaction := binary.BigEndian.Uint32(data[4:8])
	if transaction != transactionId {
		err = fmt.Errorf("transaction id should be %d, but is %d", transactionId, transaction)
		return
	}
	connectionId = data[8:]
	return
}

func sendUDP(hostAddr string, data []byte) (response []byte, err error) {
	conn, err := net.Dial(udpNet, hostAddr)
	if err != nil {
		err = fmt.Errorf("unable to connect to host %s: %s", hostAddr, err)
		return
	}
	conn.SetReadDeadline(time.Now().Add(readDeadline))
	written, err := conn.Write(data)
	if err != nil {
		err = fmt.Errorf("unable to send data to host %s : %s", hostAddr, err)
		return
	}
	if written != len(data) {
		err = fmt.Errorf("len of data written to connection does not match data passed in function: %d vs %d", written, len(data))
		return
	}
	response = make([]byte, buffSize)
	read, err := conn.Read(response)
	if err != nil {
		return
	}
	response = response[:read]
	return
}
