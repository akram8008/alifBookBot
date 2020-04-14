package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func newMessage (update tgbotapi.Update,bot *tgbotapi.BotAPI, db *sql.DB) {
	user,err := getInfoUser(update,bot,db)
	if err!=nil {return}
	log.Println("New message from: ", user, " Message: ",update.Message.Text)

	if user.Role == betypes.TextAdminRole {
		adminFunc(update.Message.Text, bot, db, user)
	}else if user.Role == betypes.TextUserRole{
		authorized (update.Message.Text,bot,db,user)
	}else {
		notAuthorized(update,bot,db,user)
	}
}

func newQuery(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
	 admin,err := getInfoUser(update,bot,db)
 	 if err!=nil {
 	 	return}
	 log.Println("New command from: ", admin, " Data:",update.CallbackQuery.Data)
	 if admin.Role!=betypes.AdminRole {
	 	return
	 }

	if update.CallbackQuery.Data == betypes.TextAddingUserQueryYes {
	   user := betypes.User{
		   ChatId:    int64(update.CallbackQuery.Message.Contact.UserID),
		   FirstName: update.CallbackQuery.Message.Contact.FirstName,
		   Phone:     update.CallbackQuery.Message.Contact.PhoneNumber,
		   Role:      "user",
	   }
	   err := dataBase.UpdateUser(db,user)
	   if err!=nil {
	   	sendErrorMessage(bot,admin.ChatId)
	   	return
	   }

		msg := tgbotapi.NewMessage(user.ChatId, betypes.TextAccepted)
		msg.ReplyMarkup = betypes.LibraryButton
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}

		msg = tgbotapi.NewMessage(admin.ChatId, fmt.Sprintf(" %s successfully added!", user.FirstName))
		msg.ReplyMarkup = betypes.LibraryButton
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}
	}else if update.CallbackQuery.Data == betypes.TextAddingUserQueryNo {
		user := betypes.User{
			ChatId:    int64(update.CallbackQuery.Message.Contact.UserID),
			FirstName: update.CallbackQuery.Message.Contact.FirstName,
			Phone:     update.CallbackQuery.Message.Contact.PhoneNumber,
			Role:      "",
		}
		err := dataBase.UpdateUser(db,user)
		if err!=nil {
			sendErrorMessage(bot,admin.ChatId)
			return
		}

		msg := tgbotapi.NewMessage(user.ChatId, betypes.TextRejected)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}

		msg = tgbotapi.NewMessage(admin.ChatId, fmt.Sprintf(" %s successfully rejected from using library!", user.FirstName))
		msg.ReplyMarkup = betypes.LibraryButton
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}
	}
}

func notAuthorized (update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB, user betypes.User) {
	var msg tgbotapi.MessageConfig

	if update.Message.Contact != nil {
		if update.Message.Contact.UserID == update.Message.From.ID {
			log.Print("Sending user's contact to admin for register")
			msgToAdmin1 := tgbotapi.NewMessage(betypes.AdminChatId, betypes.TextWantRegistration)
			if _,err := bot.Send(msgToAdmin1); err!=nil {
				log.Println("Can not send message to user ",user.FirstName," error: ",err)
				sendErrorMessage(bot,user.ChatId)
				return
			}

			msgToAdmin := tgbotapi.NewContact(betypes.AdminChatId, update.Message.Contact.PhoneNumber, update.Message.Contact.FirstName)
			msgToAdmin.ReplyMarkup = betypes.AddUserMenu
			if _,err := bot.Send(msgToAdmin); err!=nil {
				log.Println("Can not send message to user ",user.FirstName," error: ",err)
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
		log.Println("Can not send message to user ",user.FirstName," error: ",err)
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
		log.Println("Can not connect to server by the error",err)
		sendErrorMessage (bot,update.Message.Chat.ID)
		return betypes.User{},err
	}
	return user,nil
}

func authorized (message string, bot *tgbotapi.BotAPI, db *sql.DB, user betypes.User) {
	if message == betypes.TextWantLibrary{
		token := "NEWUSER"/*,err := makeToken (string(update.Message.From.ID),update.Message.From.FirstName,update.Message.From.UserName,"user")*/
		if true/*err==nil*/ {
			msg := tgbotapi.NewMessage(user.ChatId, fmt.Sprintf(betypes.TextGiveToken,token))
			log.Println(token)
			if _,err := bot.Send(msg); err!=nil {
				log.Println("Can not send message to admin error: ",err)
				return
			}
		}else {
			sendErrorMessage(bot,user.ChatId)
		}
	}else {
		msg := tgbotapi.NewMessage(user.ChatId, betypes.TextWelcome)
		msg.ReplyMarkup = betypes.TextContactSend
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}
	}
}

func adminFunc (message string, bot *tgbotapi.BotAPI, db *sql.DB, admin betypes.User) {
	if message == betypes.TextWantLibrary{
		token := "NEWADMIN"/*makeToken (string(update.Message.From.ID),update.Message.From.FirstName,update.Message.From.UserName,"admin")*/
		if true /*err == nil*/ {
			msg := tgbotapi.NewMessage(admin.ChatId, fmt.Sprintf(betypes.TextGiveToken,token))
			log.Println(token)
			if _,err := bot.Send(msg); err!=nil {
				log.Println("Can not send message to admin error: ",err)
				return
			}
		}else {
			sendErrorMessage(bot,admin.ChatId)
		}
	}else {
		msg := tgbotapi.NewMessage(admin.ChatId, betypes.TextWelcome)
		msg.ReplyMarkup = betypes.TextContactSend
		if _,err := bot.Send(msg); err!=nil {
			log.Println("Can not send message to admin error: ",err)
			return
		}
	}
}

func sendErrorMessage (bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, betypes.TextServerNotResponse)
	_, err := bot.Send(msg)
	if err!=nil {
		log.Println("Can not send error message ")
	}
}



