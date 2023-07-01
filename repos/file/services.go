package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Services struct {
}

func NewServices() *Services {
	return &Services{}
}

func (r *Services) Save(from int64, service, login, password string) error {
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

func (r *Services) Read(from int64, service string) (string, string, error) {
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

func (r *Services) Delete(from int64, service string) error {
	filename := fmt.Sprintf("data/user_%d_%s.txt", from, service)

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func (r *Services) GetAll(from int64) ([]string, error) {
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
