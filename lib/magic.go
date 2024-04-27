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
	PrintSearchUri string `json:"prints_search_uri"`
	CardFaces      []struct {
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

type Prints struct {
	Object     string `json:"object"`
	TotalCards int    `json:"total_cards"`
	HasMore    bool   `json:"has_more"`
	NextPage   string `json:"next_page"`
	Data       []struct {
		ID        string `json:"id"`
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
			Name      string `json:"name"`
			ImageUris struct {
				Small      string `json:"small"`
				Normal     string `json:"normal"`
				Large      string `json:"large"`
				Png        string `json:"png"`
				ArtCrop    string `json:"art_crop"`
				BorderCrop string `json:"border_crop"`
			} `json:"image_uris"`
		} `json:"card_faces"`
		SetName string `json:"set_name"`
	} `json:"data"`
}

func Split(r rune) bool {
	return r == '(' || r == ')'
}

func SetDive(setname string, apiuri string) string {
	imagetoreturn := ""

	response, err := http.Get(apiuri)
	if err != nil {
		fmt.Print(err.Error())
		return "error"
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		return "error parsing set details"
	}

	var prints Prints
	json.Unmarshal(responseData, &prints)

	for _, element := range prints.Data {
		if strings.EqualFold(setname, element.SetName) { //make fuzzy
			if element.Layout == "transform" {
				cardString := ""
				for _, face := range element.CardFaces {
					cardString += face.ImageUris.Large + " "
				}
				imagetoreturn = cardString
			} else {
				imagetoreturn = element.ImageUris.Large
			}
		}
	}

	if prints.HasMore {
		imagetoreturn = SetDive(setname, prints.NextPage)
	}

	if imagetoreturn == "" {
		return fmt.Sprintf("Card does not exist for set %s", setname)
	} else {
		return imagetoreturn
	}
}

func GetCard(cardname []string) string {
	commandValue := strings.Join(cardname[1:], " ")
	setSplit := strings.FieldsFunc(commandValue, Split)
	name := setSplit[0]

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

	var card Card
	json.Unmarshal(responseData, &card)

	//specific set requested. Time to search
	if len(setSplit) > 1 {
		fmt.Println(setSplit[1])
		imagetoreturn := SetDive(setSplit[1], card.PrintSearchUri)
		return imagetoreturn
	}

	if card.Layout == "transform" {
		cardString := ""
		for _, element := range card.CardFaces {
			cardString += element.ImageUris.Large + " "
		}
		return cardString
	} else {
		return card.ImageUris.Large
	}
}
