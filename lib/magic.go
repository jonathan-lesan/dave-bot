package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Card struct {
	Name      string `json:"name"`
	Layout    string `json:"layout"`
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

func GetCard(cardname []string) string {
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

	var responseObject Card

	json.Unmarshal(responseData, &responseObject)

	if responseObject.Layout == "transform" {
		cardString := ""
		for _, element := range responseObject.CardFaces {
			cardString += element.ImageUris.Large + " "
		}
		return cardString
	} else {
		return responseObject.ImageUris.Large
	}
}
