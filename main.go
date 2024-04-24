package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func gaslight() string {
	array_name := []string{
		"I never said that.",
		"I did that because I love you.",
		"I don't know why you're making such a huge deal of this.",
		"You're being overly sensitive.",
		"You are being dramatic.",
		"You are the issue, not me.",
		"You are crazy.",
		"You're being delusional.",
		"You are just insecure.",
		"You are so selfish if you don't do this for me.",
		"You're imagining things.",
		"You made me do that.",
		"You don't really feel that way.",
		"That never happened.",
		"It's not that big a deal.",
		"You're just being paranoid.",
	}

	return array_name[rand.Intn(len(array_name))]
}

//func card(cardname string) string {
//	response, err := http.Get("https://api.scryfall.com/cards/named?fuzzy=aust+com")
//}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	switch args[0] {
	case "!gaslight":
		s.ChannelMessageSend(m.ChannelID, gaslight())
	case "!card":
		s.ChannelMessageSend(m.ChannelID, "card(args[1])")
	default:
		return
	}
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
