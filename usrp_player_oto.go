// usrp_player_oto.go
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hajimehoshi/oto"
)

const (
	// USRP packet header is 32 bytes.
	headerSize = 32

	// Voice payload: 160 samples, each a 16-bit (2-byte) value = 320 bytes.
	numSamples         = 160
	sampleSize         = 2
	usrpVoiceFrameSize = numSamples * sampleSize

	// USRP type definitions
	USRP_TYPE_VOICE = 0
	// USRP_TYPE_DTMF = 1
	// USRP_TYPE_TEXT  = 2
)

// USRPHeader represents the USRP UDP packet header.
// All multi-byte numbers are transmitted in network byte order.
type USRPHeader struct {
	Eye       [4]byte // Should be "USRP"
	Seq       uint32  // Sequence counter
	Memory    uint32  // Memory identifier (or zero)
	Keyup     uint32  // Push-to-talk flag/state
	Talkgroup uint32  // Talkgroup identifier
	Type      uint32  // Payload type (0 = voice)
	Mpxid     uint32  // Reserved for future use
	Reserved  uint32  // Reserved for future use
}

func main() {
	// Command-line flag for the UDP port.
	portFlag := flag.Int("port", 1234, "UDP port to listen on for USRP packets")
	flag.Parse()

	// Set up the UDP listener.
	udpAddr := net.UDPAddr{
		Port: *portFlag,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on UDP port %d: %v\n", *portFlag, err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Printf("Listening on UDP port %d for USRP packets...", *portFlag)

	// Initialize Oto. We specify:
	//  - Sample rate: 8000 Hz
	//  - Number of channels: 1 (mono)
	//  - Bytes per sample: 2 (16-bit PCM)
	//  - A buffer size in bytes (e.g., 4096)
	ctx, err := oto.NewContext(8000, 1, 2, 4096)
	if err != nil {
		log.Fatalf("Error initializing Oto: %v", err)
	}

	player := ctx.NewPlayer()
	defer player.Close()

	// Buffer for receiving UDP packets.
	packetBuf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(packetBuf)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		if n < headerSize {
			log.Printf("Packet from %v too short: %d bytes", addr, n)
			continue
		}

		// Parse the header.
		var hdr USRPHeader
		r := bytes.NewReader(packetBuf[:headerSize])
		if err := binary.Read(r, binary.BigEndian, &hdr); err != nil {
			log.Printf("Error reading header from %v: %v", addr, err)
			continue
		}

		// Validate the signature.
		if string(hdr.Eye[:]) != "USRP" {
			log.Printf("Invalid packet signature from %v", addr)
			continue
		}

		// Process only voice packets.
		if hdr.Type != USRP_TYPE_VOICE {
			// Handle other types (e.g., DTMF, text) as needed.
			continue
		}

		// Ensure the payload is complete.
		payloadLen := n - headerSize
		if payloadLen < usrpVoiceFrameSize {
			log.Printf("Incomplete voice payload from %v: %d bytes", addr, payloadLen)
			continue
		}

		// Extract the voice payload.
		voiceData := packetBuf[headerSize : headerSize+usrpVoiceFrameSize]

		// NOTE: If you experience audio issues, you might need to convert the
		// endianness of the samples depending on how they were transmitted.
		// For example, if the samples are in big-endian order but your audio
		// system expects little-endian, you’d need to swap the byte order here.

		// Write the raw PCM data to the Oto player.
		// Oto expects a stream of bytes.
		_, err = player.Write(voiceData)
		if err != nil {
			log.Printf("Error writing audio data: %v", err)
			continue
		}
	}
}
