package taunt

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// not sent by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// not a number
	tauntID, err := strconv.Atoi(m.Content)
	if err != nil {
		return
	}

	// AoE2 has 42 unique taunts at time of writing
	if tauntID < 1 || tauntID > 42 {
		return
	}

	channelID := getChannelID(s, m.Message)
	if channelID != "" {
		play(tauntID, m.GuildID, channelID)
	}
}

func getChannelID(s *discordgo.Session, m *discordgo.Message) string {
	// Find the guild the message was sent in.
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		log.Printf("could not find guild: %s", m.GuildID)
		// Could not find guild.
		return ""
	}

	// Look for the message sender in that guild's current voice states.
	log.Println(g.VoiceStates)

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			return vs.ChannelID
		}
	}

	//log.Printf("could not find user in voice: %s", m.Author.ID)

	return ""
}
