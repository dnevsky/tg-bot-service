Бот реализован с использованием библиотеки для работы с API telegram - `github.com/go-telegram-bot-api/telegram-bot-api`

Данные хранятся в директории `data`

Хотел органичиться использованием только Dockerfile, но для сохранения данных в случае перезапуска приложения вынужнен использовать docker-compose с volume-переменной.

Так же собранный образ приложения хранится на hub.docker.com - `docker pull dnevsky/tg-bot-service`

Вывести все команды - `/help`

`/set <service> <login> <password>` - задать сервису логин и пароль.

`/get <service>` - получить логин и пароль от сервиса.

`/getall` - получить список всех сервисов.

`/del <service>` - удалить сервис.


`Makefile`:

`make build` - собрать docker образ приложения.

`make run` - запустить приложение через docker-compose.

`make shutdown` - остановить выполнение приложения.