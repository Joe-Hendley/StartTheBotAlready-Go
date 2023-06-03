package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Joe-Hendley/StartTheBotAlready-Go/internal/bot/handlers/taunt"
	"github.com/Joe-Hendley/StartTheBotAlready-Go/internal/voice"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	session, err := discordgo.New("Bot " + os.Getenv("token"))

	if err != nil {
		log.Fatalf("error creating session: %v", err)
	}

	player := voice.NewPlayer(session)
	taunt.RequestChan = player.RequestChan()

	session.AddHandler(taunt.Handler)

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	log.Println("bot is running. Press CTL-C to exit.")
	defer session.Close()

	waitForInterrupt()

	log.Println("bot shutting down")

	player.Close()

	log.Println("bot shut down successfully")
}

func waitForInterrupt() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done
}
