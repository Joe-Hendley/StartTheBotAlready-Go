package voice

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Buffer [][]byte

type Request struct {
	buffer    *Buffer
	guildID   string
	channelID string
}

func NewRequest(buffer *Buffer, guildID, channelID string) Request {
	return Request{
		buffer:    buffer,
		guildID:   guildID,
		channelID: channelID,
	}
}

// ripped from here
// https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go
func LoadSound(path string) (*Buffer, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Println("error opening dca file :", err)
		return nil, err
	}

	var opuslen int16
	var buffer Buffer

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return &buffer, nil
		}

		if err != nil {
			log.Println("error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			log.Println("error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func playSound(vc *discordgo.VoiceConnection, buffer Buffer) {

	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)
}
