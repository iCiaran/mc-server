package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
)

type String struct {
	length int32
	data   string
}

type HandshakeRequest struct {
	protocolVersion int32
	serverAddress   String
	serverPort      uint16
	nextState       int32
}

func parseVarInt(reader io.Reader) (int32, int, error) {
	var value int
	position := 0

	for {
		var current byte
		err := binary.Read(reader, binary.BigEndian, &current)
		if err != nil {
			return 0, position + 1, err
		}

		value |= int(current&0x7f) << (position * 7)
		position++

		if position > 5 {
			return 0, position, errors.New("VarInt is too big")
		}

		if (current & 0x80) == 0 {
			break
		}
	}

	return int32(value), position, nil
}

func parseString(reader io.Reader) (String, error) {
	stringLength, _, err := parseVarInt(reader)
	if err != nil {
		return String{}, err
	}

	textBytes := make([]byte, stringLength)
	_, err = reader.Read(textBytes)
	if err != nil {
		return String{}, err
	}

	return String{data: string(textBytes), length: stringLength}, nil

}

func handleConnection(conn net.Conn) {
	log.Println("Received new connection")
	defer func(conn net.Conn) {
		log.Println("Closing connection")
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection")
		}
	}(conn)

	packetLength, _, err := parseVarInt(conn)
	if err != nil {
		log.Printf("Error parsing varint: %v", err)
		return
	}
	log.Printf("Packet length is %d bytes\n", packetLength)

	packetId, n, err := parseVarInt(conn)
	if err != nil {
		log.Printf("Error parsing varint: %v", err)
		return
	}
	log.Printf("Packet id is %d\n", packetId)

	if packetId != 0 {
		return
	}

	data := make([]byte, packetLength-int32(n))
	err = binary.Read(conn, binary.BigEndian, &data)
	if err != nil {
		log.Printf("Error reading data: %v", err)
		return
	}
	dataReader := bytes.NewReader(data)

	protocolVersion, _, err := parseVarInt(dataReader)
	if err != nil {
		log.Printf("Error parsing protocol version: %v", err)
		return
	}
	log.Printf("Protocol version is %d\n", protocolVersion)

	serverAddress, err := parseString(dataReader)
	if err != nil {
		log.Printf("Error parsing server address: %v", err)
		return
	}
	log.Printf("Server address is %s\n", serverAddress.data)

	serverPortBytes := make([]byte, 2)
	_, err = dataReader.Read(serverPortBytes)
	if err != nil {
		log.Printf("Error parsing server port: %v", err)
		return
	}

	serverPort := binary.BigEndian.Uint16(serverPortBytes)
	log.Printf("Server port is %d\n", serverPort)

	nextState, _, err := parseVarInt(conn)
	if err != nil {
		log.Printf("Error parsing nextstate: %v", err)
		return
	}
	log.Printf("Next state is %d\n", nextState)

	handShakePacket := HandshakeRequest{protocolVersion: protocolVersion, serverAddress: serverAddress, serverPort: serverPort, nextState: nextState}
	log.Printf("Handshake packet is %+v\n", handShakePacket)

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
