package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func newMessage (update tgbotapi.Update,bot *tgbotapi.BotAPI, db *sql.DB) {
	user,err := getInfoUser(update,bot,db)
	if err!=nil {return}
	log.Println("New message from: ", user)

	if user.Role == "admin" {
		//adminFunc(update, bot, db, user)
	}else if user.Role == "user"{
		//authorized (update,bot,db, user)
	}else {
		notAuthorized(update,bot,db,user)
	}
}




func notAuthorized (update tgbotapi.Update,bot *tgbotapi.BotAPI, db *sql.DB, user betypes.User) {
	var msg tgbotapi.MessageConfig

	if update.Message.Contact != nil {
		if update.Message.Contact.UserID == update.Message.From.ID {
			log.Print("Sending user's contact to admin for register")
			msgToAdmin1 := tgbotapi.NewMessage(betypes.AdminChatId, betypes.TextWantRegistration)
			if _,err := bot.Send(msgToAdmin1); err!=nil {
				sendErrorMessage(bot,user.ChatId)
				return
			}

			msgToAdmin := tgbotapi.NewContact(betypes.AdminChatId, update.Message.Contact.PhoneNumber, update.Message.Contact.FirstName)
			msgToAdmin.ReplyMarkup = betypes.AddUserMenu
			if _,err := bot.Send(msgToAdmin); err!=nil {
				sendErrorMessage(bot,user.ChatId)
				return
			}

			user.Role = betypes.StatusWait
			updateUser(user,bot,db)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, betypes.TextRegistrationSent)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	}else if user.Role == betypes.StatusWait{
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, betypes.TextRegistrationWait)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, betypes.TextWelcome)
		msg.ReplyMarkup = betypes.ContactButton
	}

	if _,err := bot.Send(msg); err!=nil {
		sendErrorMessage(bot,user.ChatId)
		return
	}
}

func updateUser (user betypes.User, bot *tgbotapi.BotAPI,db *sql.DB) {
	err := dataBase.UpdateUser(db,user)
	if err != nil {
		log.Println("Can not update information of user ")
		sendErrorMessage (bot,betypes.AdminChatId)
		return
	}
}

func getInfoUser (update tgbotapi.Update,bot *tgbotapi.BotAPI, db *sql.DB) (betypes.User, error) {
	user := betypes.User{ChatId:update.Message.Chat.ID}
	err := dataBase.InfoUserDB(db,&user)
	if err != nil {
		log.Println("Can not connect to server ")
		sendErrorMessage (bot,update.Message.Chat.ID)
		return betypes.User{},err
	}
	return user,nil
}


/*
func authorized (update tgbotapi.Update,bot *tgbotapi.BotAPI) {
	if update.Message.Text == betypes.TextWantLibrary{
		token,err := makeToken (string(update.Message.From.ID),update.Message.From.FirstName,update.Message.From.UserName,"user")
		if err==nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваш кабинет открыт на 10 часов на ползование.\n\n https://alif-library/?data=%s",token))
			bot.Send(msg)
		}else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер временно не доступно! Повторите попозже!")
			bot.Send(msg)
		}
	}else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в электронная библиотека Алифа!!")
		msg.ReplyMarkup = libraryButton
		bot.Send(msg)
	}
}
*/
/*
func adminFunc (update tgbotapi.Update,bot *tgbotapi.BotAPI) {
	if update.Message.Text == "Хочу в библиотеку 🏠"{
		token,err := makeToken (string(update.Message.From.ID),update.Message.From.FirstName,update.Message.From.UserName,"admin")
		if err == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваш кабинет открыт на 10 часов на ползование.\n\n https://alif-library/?data=%s",token))
			log.Println(token)
			bot.Send(msg)
		}else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сервер временно не доступно! Повторите попозже!")
			bot.Send(msg)
		}
	}else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в электронная библиотека Алифа!!")
		msg.ReplyMarkup = libraryButton
		bot.Send(msg)
	}
}
*/



func sendErrorMessage (bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Серевер не доступенно! повторите попоезже!")
	_, err := bot.Send(msg)
	if err!=nil {
		log.Println("Can not send error message to user about Database connection problems")
	}
}




