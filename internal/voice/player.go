package voice

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	session *discordgo.Session
	guilds  map[string]*guild

	requestChan chan Request
	doneChan    chan struct{}
}

func NewPlayer(s *discordgo.Session) *Player {
	p := &Player{
		session: s,
		guilds:  map[string]*guild{},

		requestChan: make(chan Request),
		doneChan:    make(chan struct{}),
	}

	readyChan := make(chan struct{})
	log.Println("starting player")
	go p.process(readyChan)

	select {
	case <-readyChan:
		return p
	case <-time.After(time.Second):
		log.Fatal("player not started after 1 second")
	}

	return p
}

func (p *Player) RequestChan() chan<- Request {
	return p.requestChan
}

func (p *Player) Close() {
	log.Println("closing player")
	close(p.requestChan)

	wg := sync.WaitGroup{}

	for _, g := range p.guilds {
		wg.Add(1)
		go func(g *guild) {
			g.close()
			wg.Done()
		}(g)
	}

	wg.Wait()
	log.Println("guilds closed")

	select {
	case <-p.doneChan:
		return
	case <-time.After(10 * time.Second):
		log.Fatal("audio player failed to close after 10 seconds")
	}

}

func (p *Player) process(readyChan chan struct{}) {
	log.Println("audio player started")
	readyChan <- struct{}{}

	for r := range p.requestChan {
		log.Printf("player req: %v\n", r)
		var g *guild
		g, ok := p.guilds[r.guildID]
		if !ok {
			g = &guild{session: p.session, ID: r.guildID, requestChan: make(chan Request, maxRequestsPerGuild), doneChan: make(chan struct{})}
			p.guilds[r.guildID] = g
			go g.process()
		}

		g.requestChan <- r
	}
	log.Println("audio player shutting down")

	p.doneChan <- struct{}{}
}
