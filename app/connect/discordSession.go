package connect

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var DG *discordgo.Session

func CreateDiscordSession() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Fail to load env")
	}
	token := os.Getenv("Discord_Token")

	DG, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Fail to create bot Session")
	}
}
