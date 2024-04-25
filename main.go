package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Card struct {
	Name      string `json:"name"`
	ImageUris struct {
		Small      string `json:"small"`
		Normal     string `json:"normal"`
		Large      string `json:"large"`
		Png        string `json:"png"`
		ArtCrop    string `json:"art_crop"`
		BorderCrop string `json:"border_crop"`
	} `json:"image_uris"`
	CardFaces []struct {
		Object         string   `json:"object"`
		Name           string   `json:"name"`
		ManaCost       string   `json:"mana_cost"`
		TypeLine       string   `json:"type_line"`
		OracleText     string   `json:"oracle_text"`
		Colors         []string `json:"colors"`
		Loyalty        string   `json:"loyalty"`
		Artist         string   `json:"artist"`
		ArtistID       string   `json:"artist_id"`
		IllustrationID string   `json:"illustration_id"`
		ImageUris      struct {
			Small      string `json:"small"`
			Normal     string `json:"normal"`
			Large      string `json:"large"`
			Png        string `json:"png"`
			ArtCrop    string `json:"art_crop"`
			BorderCrop string `json:"border_crop"`
		} `json:"image_uris"`
	} `json:"card_faces"`
}

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

func card(cardname []string) string {
	name := strings.Join(cardname[1:], " ")
	response, err := http.Get(fmt.Sprintf("https://api.scryfall.com/cards/named?fuzzy=%s", name))

	if err != nil {
		fmt.Print(err.Error())
		return "error"
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		return "error parsing response"
	}
	//fmt.Println(string(responseData))

	var responseObject Card

	json.Unmarshal(responseData, &responseObject)

	if responseObject.CardFaces != nil {
		cardString := ""
		for _, element := range responseObject.CardFaces {
			cardString += element.ImageUris.Large + " "
		}
		return cardString
	} else {
		return responseObject.ImageUris.Large
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	switch args[0] {
	case "!gaslight":
		s.ChannelMessageSend(m.ChannelID, gaslight())
	case "!card":
		s.ChannelMessageSend(m.ChannelID, card(args))
	default:
		return
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

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
