package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"
)

// player struct
type player struct {
	address    *net.UDPAddr
	id         int
	scene      int
	posx       float32
	posy       float32
	lastUpdate int64
}

// players is a map of all players currently connected to the server
var players = make(map[string]*player)
var playerID int = 0

func receivedPosition(conn *net.UDPConn, addr *net.UDPAddr, msg []byte) {
	// Decode scene as int32
	scene := binary.LittleEndian.Uint32(msg[0:4])

	// Decode x position as float32
	xPos := binary.LittleEndian.Uint32(msg[4:8])
	xPosFloat := math.Float32frombits(xPos)

	// Decode y position as float32
	yPos := binary.LittleEndian.Uint32(msg[8:12])
	yPosFloat := math.Float32frombits(yPos)

	playerAddress := addr.String()

	// Check if player is already in the map
	if _, ok := players[playerAddress]; ok {
		// Update player
		players[playerAddress].scene = int(scene)
		players[playerAddress].posx = xPosFloat
		players[playerAddress].posy = yPosFloat
		players[playerAddress].lastUpdate = getTime()
	} else {
		// Add player to map
		players[playerAddress] = &player{
			address:    addr,
			scene:      int(scene),
			posx:       xPosFloat,
			posy:       yPosFloat,
			lastUpdate: getTime(),
			id:         playerID,
		}
		playerID++
		fmt.Println("Player", players[playerAddress].id, "connected")
	}
}

// send all players to all players
func sendAllPlayersPos(conn *net.UDPConn) {

	timeout := getTime() - 120

	for _, player := range players {

		// Check if player is still connected
		if player.lastUpdate < timeout {
			// Remove player from map
			delete(players, player.address.String())
			fmt.Println("Player", player.id, "disconnected")
			continue
		}

		// Message type
		msgType := []byte{0x02}
		// Encode player id as int32
		id := make([]byte, 4)
		binary.LittleEndian.PutUint32(id, uint32(player.id))

		// Encode scene as int32
		scene := make([]byte, 4)
		binary.LittleEndian.PutUint32(scene, uint32(player.scene))

		// Encode x position as float32
		xPos := make([]byte, 4)
		binary.LittleEndian.PutUint32(xPos, math.Float32bits(player.posx))

		// Encode y position as float32
		yPos := make([]byte, 4)
		binary.LittleEndian.PutUint32(yPos, math.Float32bits(player.posy))

		// Encode all data
		data := append(msgType, id...)
		data = append(data, scene...)
		data = append(data, xPos...)
		data = append(data, yPos...)

		// Send data to all players except the player itself
		for _, thePlayer := range players {
			if thePlayer.address.String() != player.address.String() {
				sendResponse(conn, thePlayer.address, data)
			}
		}
	}

	// Wait and call this function again
	time.Sleep(200 * time.Millisecond)
	sendAllPlayersPos(conn)
}
