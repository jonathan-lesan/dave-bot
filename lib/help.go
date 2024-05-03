package lib

func SendHelp() string {
	helpMessage := ""

	helpMessage += "Below are the available commands for the GasLightBot to use:\n\n"
	helpMessage += "!gaslight - Will send a gaslighting message back to you. Hope you're ready.\n"
	helpMessage += "!card - looks up the following text in Scryfall for a MTG card you wish to show. Add the set name/abbreviation in parenthesis for different art ex: !card lord wind grace (sld)\n"

	return helpMessage
}
