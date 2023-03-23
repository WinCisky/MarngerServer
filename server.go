package main

import (
	"fmt"
	"net"
)

func main() {
	buf := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	// start sending all players to all players
	go sendAllPlayersPos(ser)

	fmt.Println("Server started")

	for {
		n, remoteaddr, err := ser.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		// Decode the message based on the first byte
		switch buf[0] {
		case 0x01: // ping
			answer := []byte{0x01}
			go sendResponse(ser, remoteaddr, answer)
		case 0x02: // pos -> no answer
			go receivedPosition(ser, remoteaddr, buf[1:n])
		default:
			fmt.Println("Unknown message type:", buf[0])
		}
	}
}
