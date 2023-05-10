package tg

import (
	"fmt"
	"log"
	"os"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBotApi struct {
	bot *tgapi.BotAPI
}

func NewTGBotApi(token string) (*TGBotApi, error) {
	bot, err := tgapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		fmt.Printf("Не удалось инициализировать бота.\n%s", err.Error())
		return nil, err
	}

	bot.Debug = true

	return &TGBotApi{bot: bot}, nil
}

func (t *TGBotApi) StartLongPoll() {
	updateConfig := tgapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		go t.handler(update)
	}
}

func (t *TGBotApi) SendMessage(to int64, text string) tgapi.Message {
	msg := tgapi.NewMessage(to, text)

	message, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Не удалось отправить сообщение.\n%s\nrequest:%#v\nresponse:%#v\n", err.Error(), msg, message)
		return tgapi.Message{}
	}

	return message
}

func (t *TGBotApi) DeleteMessage(chat_id int64, message_id int) error {
	deleteMessageConfig := tgapi.DeleteMessageConfig{
		ChatID:    chat_id,
		MessageID: int(message_id),
	}

	_, err := t.bot.Send(deleteMessageConfig)
	if err != nil {
		return err
	}

	return nil
}
