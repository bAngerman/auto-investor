package discord

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Message using the discord bot
func Message(message string) {

	token := os.Getenv("DISCORD_TOKEN")
	user := os.Getenv("DISCORD_USER")

	if token == "" || user == "" {
		fmt.Println("Error getting env DISCORD_TOKEN, DISCORD_USER")
		return
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
	}

	s.Identify.Intents = discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages |
		discordgo.PermissionManageMessages

	defer s.Close()

	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	channel, err := s.UserChannelCreate(user)

	if err != nil {
		fmt.Println("Error getting channel:", err)
	}

	_, err = s.ChannelMessageSend(channel.ID, message)

	if err != nil {
		fmt.Println("Error sending DM:", err)
	}
}
