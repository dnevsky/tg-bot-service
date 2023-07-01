package tg

import (
	"context"
	"fmt"
	"log"

	"github.com/dnevsky/tg-bot-service/repos"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBotApi struct {
	bot   *tgapi.BotAPI
	repos *repos.Repos
}

func NewTGBotApi(token string) (*TGBotApi, error) {
	bot, err := tgapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("Не удалось инициализировать бота.\n%s", err.Error())
		return nil, err
	}

	bot.Debug = false

	return &TGBotApi{bot: bot}, nil
}

func (t *TGBotApi) StartLongPoll(ctx context.Context, repos *repos.Repos) {
	t.repos = repos
	updateConfig := tgapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case update := <-updates:
			go t.handler(ctx, update)
		case <-ctx.Done():
			log.Println("Stopping longpoll...")
			return
		}
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
