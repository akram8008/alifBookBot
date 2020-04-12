package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func main () {
	log.Print("start application")
	port := os.Getenv(betypes.EnvPort)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Println("Listening on server: ", betypes.BotWebhook + ":" + port)

	db := dataBase.Connect()
	bot := botConnect()

	start(bot,db)
}




func start (bot *tgbotapi.BotAPI,db *sql.DB) {

	updates := bot.ListenForWebhook("/")

	for update := range updates {
		if update.Message != nil {
			newMessage(update,bot,db)
		}
		if update.CallbackQuery != nil {
			//newQuery(update,db)
		}
	}
}



func botConnect () *tgbotapi.BotAPI {
	log.Println("Connecting to bot api! ")
	bot, err := tgbotapi.NewBotAPI(betypes.BotToken)
	if err != nil {
		log.Fatal("Can't connect to bot api")
	}
	log.Printf("Authorized bot api - %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(betypes.BotWebhook))
	if err != nil {
		log.Fatal("Can't connect set webhook of telegram-bot",err)
	}
	log.Println("Webhook set")
	return bot
}

