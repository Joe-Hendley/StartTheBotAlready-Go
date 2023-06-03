package voice

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

const maxRequestsPerGuild = 16

type guild struct {
	session *discordgo.Session
	ID      string

	requestChan chan Request
	doneChan    chan struct{}
	done        bool
}

func (g *guild) process() {
	var vc *discordgo.VoiceConnection

	log.Printf("guild %s started", g.ID)
	for r := range g.requestChan {
		log.Printf("guild req: %v\n", r)
		if g.done { // we don't want to continue playing
			continue
		}

		switch {
		case vc != nil && vc.ChannelID == r.channelID: // already connected and in the right channel
			log.Println("vc & channel")
			playSound(vc, *r.buffer)

		case vc != nil: // connected to the wrong channel
			log.Println("vc")
			vc.Disconnect()
			var err error
			vc, err = g.session.ChannelVoiceJoin(g.ID, r.channelID, false, true)
			if err != nil {
				log.Printf("cannot connect to voice channel: %v", err)
				continue
			}
			playSound(vc, *r.buffer)

		default: // not connected
			log.Println("neither")
			var err error
			vc, err = g.session.ChannelVoiceJoin(g.ID, r.channelID, false, true)
			if err != nil {
				log.Printf("cannot connect to voice channel: %v", err)
				continue
			}
			playSound(vc, *r.buffer)
		}

		if len(g.requestChan) == 0 && vc != nil {
			log.Printf("disconnecting from guild %s\n", g.ID)
			vc.Disconnect()
			vc = nil
		}
	}
	g.doneChan <- struct{}{}
}

func (g *guild) close() {
	log.Printf("closing guild: %s\n", g.ID)
	close(g.requestChan)
	g.done = true

	select {
	case <-g.doneChan:
		return
	case <-time.After(10 * time.Second): // long timeout for long audio
		log.Fatalf("guild player[%s] failed to close after 10 seconds", g.ID)
	}
}
