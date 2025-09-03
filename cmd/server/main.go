package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/iCiaran/mc-server/internal/packets"
)

func decodePacket(reader io.Reader, state packets.VarInt) (interface{}, error) {
	packetLength, _, err := packets.DeserializeVarInt(reader)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, packetLength)
	_, err = reader.Read(buf)
	if err != nil {
		return nil, err
	}

	packetReader := bytes.NewReader(buf)
	packetId, _, err := packets.DeserializeVarInt(packetReader)
	if err != nil {
		return nil, err
	}

	if state == 0 && packetId == 0x00 {
		log.Println("Deserializing Intention")
		intention, _, err := packets.DeserializeIntention(packetReader)
		return intention, err
	} else if state == 1 && packetId == 0x00 {
		log.Println("Deserializing StatusRequest")
		statusRequest, _, err := packets.DeserializeStatusRequest(packetReader)
		return statusRequest, err
	} else if state == 1 && packetId == 0x01 {
		log.Println("Deserializing PingRequest")
		pingRequest, _, err := packets.DeserializePingRequest(packetReader)
		return pingRequest, err
	}

	return nil, fmt.Errorf("unknown packet (state: %d, id: %x)", state, packetId)
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v\n", err)
		}
	}()

	state := packets.VarInt(0)

	intention, err := decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding intent: %v\n", err)
		return
	}

	state = intention.(packets.Intention).Intent

	_, err = decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding statusRequest: %v\n", err)
		return
	}

	statusResponse, err := packets.StatusResponse{
		Response: packets.StatusResponseJson{
			Version: packets.StatusResponseVersion{
				Name:     "1.21",
				Protocol: 767,
			},
			Players: packets.StatusResponsePlayers{
				Max:    10,
				Online: 0,
			},
			Description: packets.StatusResponseDescription{
				Text: "Ciaran woz ere",
			},
			EnforceSecureChat: false,
		},
	}.Serialize()
	if err != nil {
		log.Printf("Error serializing statusResponse: %v\n", err)
		return
	}

	_, err = conn.Write(statusResponse)
	if err != nil {
		log.Printf("Error writing statusResponse: %v\n", err)
		return
	}

	pingRequest, err := decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding pingRequest: %v\n", err)
		return
	}

	pongResponse, err := packets.PongResponse{
		Timestamp: pingRequest.(packets.PingRequest).Timestamp,
	}.Serialize()
	if err != nil {
		log.Printf("Error serializing pongResponse: %v\n", err)
		return
	}

	_, err = conn.Write(pongResponse)
	if err != nil {
		log.Printf("Error writing pongResponse: %v\n", err)
	}

}

func main() {
	ln, err := net.Listen("tcp4", ":25565")
	if err != nil {
		log.Printf("Error listening on tcp: %v\n", err)
		os.Exit(1)
	}

	log.Println("Listening on tcp:", ln.Addr())

	for {
		log.Println("Accepting connection")
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)
	}
}
