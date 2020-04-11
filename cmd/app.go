package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func newMessage (update tgbotapi.Update,bot *tgbotapi.BotAPI, db *sql.DB) {
	log.Println("New message from: ", update.Message.From.FirstName)

	user := betypes.User{ChatId:update.Message.Chat.ID}
	user,ok,err := dataBase.IsUserExist(db,user)
	if err!=nil {
		log.Println("Can not connect to server ")
		sendErrorMessage (bot,update.Message.Chat.ID)
		return
	}

	if user.Role == "admin" {
		adminFunc(update, bot)
	}else if user.Role == "user"{
		authorized (update,bot)
	}else {
		notAuthorized(&update,bot)
	}
}



func sendErrorMessage (bot *tgbotapi.BotAPI, chatId int64) {
	msg = tgbotapi.NewMessage(chatId, "Серевер не доступенно! повторите попоезже!")
}