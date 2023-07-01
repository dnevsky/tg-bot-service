package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dnevsky/tg-bot-service/repos"
	"github.com/dnevsky/tg-bot-service/repos/postgres"
	"github.com/dnevsky/tg-bot-service/tg"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	godotenv.Load(".env")

	if err := initConfig(); err != nil {
		log.Fatalf("error init configs: %s", err.Error())
	}

	bot, err := tg.NewTGBotApi(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		log.Print("Shutdown bot...")
		cancel()
	}()

	log.Println("Start longpoll...")

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to init db: %s", err.Error())
	}

	repos := repos.NewRepos(db)
	bot.StartLongPoll(ctx, repos)

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("default")
	return viper.ReadInConfig()
}
