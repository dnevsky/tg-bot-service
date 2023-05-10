package tg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	dur = time.Minute
)

func (t *TGBotApi) handler(update tgapi.Update) {
	// даём знать, что мы работаем ток с сообщениями
	if update.Message == nil {
		return
	}

	cmd := update.Message.Text

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

		err := saveService(update.Message.From.ID, service, login, password)
		if err != nil {
			t.SendMessage(update.Message.From.ID, err.Error())
		}

		msg := t.SendMessage(update.Message.Chat.ID, "Успешно сохранено.")

		go t.deleteMessageAfter(msg.Chat.ID, msg.MessageID, dur)
		go t.deleteMessageAfter(update.Message.Chat.ID, update.Message.MessageID, dur)

	case strings.HasPrefix(cmd, "/getall"):
		services, err := getAll(update.Message.From.ID)
		if err != nil {

			t.SendMessage(update.Message.From.ID, err.Error())
			return
		}

		msg := "list of all services:\n\n"

		if len(services) == 0 {
			msg = msg + "not found"
			t.SendMessage(update.Message.From.ID, msg)
			return
		}

		for _, service := range services {
			login, _, _ := readService(update.Message.From.ID, service)
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

		login, password, err := readService(update.Message.From.ID, words[1])
		if err != nil {
			var msgText string
			if errors.Is(err, os.ErrNotExist) {
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

		err := deleteService(update.Message.From.ID, words[1])
		if err != nil {
			var msgText string
			if errors.Is(err, os.ErrNotExist) {
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

func saveService(from int64, service, login, password string) error {
	filename := fmt.Sprintf("data/user_%d_%s.txt", from, service)

	file, err := os.Create(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Не удалось создать новый файл.\n%s", err.Error()))
	}
	defer file.Close()

	content := fmt.Sprintf("ID: %d\nService: %s\nUsername: %s\nPassword: %s",
		from, service, login, password)

	if err := ioutil.WriteFile(filename, []byte(content), 0600); err != nil {
		return errors.New(fmt.Sprintf("Не удалось записать данные в файл.\n%s", err.Error()))
	}

	return nil
}

func readService(from int64, service string) (string, string, error) {
	filename := fmt.Sprintf("data/user_%d_%s.txt", from, service)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(content), "\n")

	var login, password string
	for _, line := range lines {
		pair := strings.SplitN(line, ": ", 2)

		if len(pair) != 2 {
			continue
		}

		switch pair[0] {
		case "Username":
			login = pair[1]
		case "Password":
			password = pair[1]
		}
	}

	return login, password, nil
}

func deleteService(from int64, service string) error {
	filename := fmt.Sprintf("data/user_%d_%s.txt", from, service)

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func getAll(from int64) ([]string, error) {
	pattern := fmt.Sprintf("user_%d_", from)

	files, err := ioutil.ReadDir("data")
	if err != nil {
		return nil, err
	}

	services := []string{}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), pattern) {
			service := strings.TrimPrefix(file.Name(), pattern)
			service = strings.TrimSuffix(service, ".txt")

			services = append(services, service)
		}
	}

	return services, nil
}

func (t *TGBotApi) deleteMessageAfter(chat_id int64, message_id int, duration time.Duration) {
	time.Sleep(duration)
	t.DeleteMessage(chat_id, message_id)
}
