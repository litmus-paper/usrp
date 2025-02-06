# USRP UDP Audio Player

This project is a simple Go terminal application that listens for USRP UDP packets and streams 16-bit PCM voice data to your PC speakers using the [Oto](https://github.com/hajimehoshi/oto) audio library. It is cross-platform and works on Windows (as well as Linux and macOS).

## Overview

The application receives UDP packets following the USRP protocol—a protocol originally designed to interface GNU Radio with Asterisk. The USRP protocol can carry voice, DTMF, and text data. This program currently processes only voice packets.

## USRP Protocol Details

USRP packets consist of a fixed 32-byte header followed by a payload. The header is defined as follows:

```c
struct _chan_usrp_bufhdr {
    char eye[4];        // Verification string. Should be "USRP"
    uint32_t seq;       // Sequence counter (network byte order)
    uint32_t memory;    // Memory identifier (or zero)
    uint32_t keyup;     // Push-to-talk (PTT) state indicator
    uint32_t talkgroup; // Trunk talk group identifier
    uint32_t type;      // Payload type:
                        //   0 = Voice data
                        //   1 = DTMF tone information
                        //   2 = Text data
    uint32_t mpxid;     // Multiplexer ID (reserved for future use)
    uint32_t reserved;  // Reserved for future use
};

## Key Points:

    Verification String (eye):
    Must be "USRP" to ensure the packet is valid.

    Sequence Number (seq):
    Increments with each packet to help detect packet loss.

    Memory (memory):
    A memory identifier (or zero if unused).

    Push-to-Talk State (keyup):
    Indicates whether the transmitter is active.

    Talkgroup (talkgroup):
    Identifies the talk group, useful in trunked radio systems.

    Payload Type (type):
        0 (Voice): The payload contains voice data.
        The voice payload is 160 samples of 16-bit linear PCM, equaling 320 bytes per packet. At an 8 kHz sampling rate, this represents approximately 20 milliseconds of audio.
        1 (DTMF): The payload contains DTMF tone information.
        2 (Text): The payload contains text data.

    Multiplexer ID (mpxid) & Reserved (reserved):
    These fields are reserved for future use.

## Requirements

    Go (version 1.XX or later recommended)
    Oto – a pure Go audio library
    A network connection for receiving UDP packets following the USRP protocol

## Building
1.  Install Go:
    Download and install Go from golang.org/dl.

2.  Download the Oto Package:
    Open your terminal (Command Prompt or PowerShell on Windows) and run:

go get github.com/hajimehoshi/oto

3.  Build the Application:
Save the source code as usrp_player_oto.go (see project source for details) and build it:

    On Windows:

go build -o usrp_player_oto.exe usrp_player_oto.go

On Linux/macOS:

        go build -o usrp_player_oto usrp_player_oto.go

## Usage

Run the compiled executable with the -port flag to specify the UDP port. For example, to listen on UDP port 1234:

./usrp_player_oto -port=1234

The program will start listening for USRP UDP packets. When a valid voice packet is received, it extracts the 320-byte payload (160 16-bit samples) and streams the audio to your PC speakers.
Windows Setup

This application uses the Oto library, which is written in pure Go and does not require any additional C libraries. Simply ensure you have Go installed and follow the building instructions above.
License

This project is provided under the terms of the GNU General Public License (GPL) as used by the original USRP channel module. Please refer to the source code for additional licensing details.
Contributing

Contributions, issues, and feature requests are welcome. Feel free to fork the repository and submit pull requests.

Enjoy streaming your USRP voice data directly to your speakers!


---

This `README.md` provides an overview of the project, details about the USRP protocol, requirements, building instructions, and usage details. You can modify and expand it as needed for your specific project.

```
