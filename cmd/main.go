package main

import (
	"log"
	"os"

	"github.com/dnevsky/tg-bot-service/tg"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	bot, err := tg.NewTGBotApi(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	bot.StartLongPoll()

}
