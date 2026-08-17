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
	} else if state == 2 && packetId == 0x00 {
		log.Println("Deserializing LoginStart")
		loginStart, _, err := packets.DeserializeLoginStart(packetReader)
		return loginStart, err
	} else if state == 2 && packetId == 0x03 {
		log.Println("Deserializing LoginAcknowledged")
		loginAcknowledged, _, err := packets.DeserializeLoginAcknowledged(packetReader)
		return loginAcknowledged, err
	}

	return nil, fmt.Errorf("unknown packet (state: %d, id: %x)", state, packetId)
}

func handleStatus(conn net.Conn, state packets.VarInt) error {
	_, err := decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding statusRequest: %v\n", err)
		return err
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
		return err
	}

	_, err = conn.Write(statusResponse)
	if err != nil {
		log.Printf("Error writing statusResponse: %v\n", err)
		return err
	}

	pingRequest, err := decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding pingRequest: %v\n", err)
		return err
	}

	pongResponse, err := packets.PongResponse{
		Timestamp: pingRequest.(packets.PingRequest).Timestamp,
	}.Serialize()
	if err != nil {
		log.Printf("Error serializing pongResponse: %v\n", err)
		return err
	}

	_, err = conn.Write(pongResponse)
	if err != nil {
		log.Printf("Error writing pongResponse: %v\n", err)
		return err
	}

	return nil
}

func handleLogin(conn net.Conn, state packets.VarInt) error {
	loginStart, err := decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding loginStart: %v\n", err)
		return err
	}

	log.Printf("LoginStart: %v\n", loginStart)

	loginFinished := packets.LoginFinished{
		Name:       loginStart.(packets.LoginStart).Name,
		UUID:       loginStart.(packets.LoginStart).UUID,
		Properties: 0,
		Strict:     false,
	}
	log.Printf("LoginFinished: %v\n", loginFinished)

	loginFinishedSerialized, err := loginFinished.Serialize()
	if err != nil {
		log.Printf("Error serializing loginFinished: %v\n", err)
		return err
	}

	_, err = conn.Write(loginFinishedSerialized)
	if err != nil {
		log.Printf("Error writing loginFinished: %v\n", err)
		return err
	}

	_, err = decodePacket(conn, state)
	if err != nil {
		log.Printf("Error decoding loginAck: %v\n", err)
		return err
	}

	return nil
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

	if state == 1 {
		err = handleStatus(conn, state)
	} else if state == 2 {
		err = handleLogin(conn, state)
	}

	if err != nil {
		log.Printf("Error handling post intent: %v\n", err)
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
