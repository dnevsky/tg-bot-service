package tg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	dur = time.Minute
)

func (t *TGBotApi) handler(ctx context.Context, update tgapi.Update) {
	select {
	case <-ctx.Done():
		log.Println("Handler canceled. Exiting...")
		return
	default:
		// работаем
	}

	// даём знать, что мы работаем ток с сообщениями
	if update.Message == nil {
		return
	}

	cmd := update.Message.Text
	log.Println(update.Message.From.ID, cmd)

	switch {
	case strings.HasPrefix(cmd, "/help"):
		msgText := "/set <service> <login> <password> - задать сервису логин и пароль\n/get <service> - получить логин и пароль по сервису\n/getall - получить список сервисов\n/del <service> - удалить логин и пароль по сервису\n/help - посмотреть все команды"

		t.SendMessage(update.Message.Chat.ID, msgText)

	case strings.HasPrefix(cmd, "/set"):
		words := strings.Fields(cmd)
		if len(words) != 4 {
			msgText := "/set <service> <login> <password>"

			t.SendMessage(update.Message.Chat.ID, msgText)
			return
		}

		service := words[1]
		login := words[2]
		password := words[3]

		err := t.repos.Save(update.Message.From.ID, service, login, password)
		if err != nil {
			t.SendMessage(update.Message.From.ID, err.Error())
			log.Println(err.Error())
			return
		}

		msg := t.SendMessage(update.Message.Chat.ID, "Успешно сохранено.")

		go t.deleteMessageAfter(msg.Chat.ID, msg.MessageID, dur)
		go t.deleteMessageAfter(update.Message.Chat.ID, update.Message.MessageID, dur)

	case strings.HasPrefix(cmd, "/getall"):
		services, err := t.repos.GetAll(update.Message.From.ID)
		if err != nil {
			t.SendMessage(update.Message.From.ID, err.Error())
			log.Println(err.Error())
			return
		}

		msg := "Список всех сервисов:\n\n"

		if len(services) == 0 {
			msg = msg + "Сервисов нет"
			t.SendMessage(update.Message.From.ID, msg)
			return
		}

		for _, service := range services {
			login, _, _ := t.repos.Read(update.Message.From.ID, service)
			msg = fmt.Sprintf("%s%s %s\n", msg, service, login)
		}

		t.SendMessage(update.Message.From.ID, msg)

	case strings.HasPrefix(cmd, "/get"):
		words := strings.Fields(cmd)
		if len(words) != 2 {
			msgText := "/get <service>"

			t.SendMessage(update.Message.Chat.ID, msgText)
			return
		}

		login, password, err := t.repos.Read(update.Message.From.ID, words[1])
		if err != nil {
			var msgText string
			if errors.Is(err, sql.ErrNoRows) {
				msgText = "Сервис не найден."
			} else {
				msgText = fmt.Sprintf("Произошла ошибка во время чтения информации.\n%s", err.Error())
			}

			t.SendMessage(update.Message.Chat.ID, msgText)

			return
		}

		msgText := fmt.Sprintf("login: %s\npassword: %s", login, password)

		msg := t.SendMessage(update.Message.Chat.ID, msgText)

		go t.deleteMessageAfter(msg.Chat.ID, msg.MessageID, dur)
		go t.deleteMessageAfter(update.Message.Chat.ID, update.Message.MessageID, dur)

	case strings.HasPrefix(cmd, "/del"):
		words := strings.Fields(cmd)
		if len(words) != 2 {
			msgText := "/del <service>"

			t.SendMessage(update.Message.Chat.ID, msgText)
			return
		}

		err := t.repos.Delete(update.Message.From.ID, words[1])
		if err != nil {
			var msgText string
			if errors.Is(err, sql.ErrNoRows) {
				msgText = "Сервис не найден."
			} else {
				msgText = fmt.Sprintf("Произошла ошибка во время удаления информации.\n%s", err.Error())
			}

			t.SendMessage(update.Message.Chat.ID, msgText)

			return
		}

		t.SendMessage(update.Message.Chat.ID, "Сервис успешно удалён.")

	}
}

func regexMatchString(pattern string, s string) (result bool) {
	ok, _ := regexp.MatchString(pattern, s)
	return ok
}

func (t *TGBotApi) deleteMessageAfter(chat_id int64, message_id int, duration time.Duration) {
	time.Sleep(duration)
	t.DeleteMessage(chat_id, message_id)
}
