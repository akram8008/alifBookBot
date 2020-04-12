package betypes

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type User struct  {
	Id        int64
	ChatId    int64
	FirstName string
	Phone     string
	Role      string
}



var libraryButton = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(TextWantLibrary)),)

var ContactButton = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButtonContact(TextContactSend)),)

var AddUserMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(TextAccept,TextQueryYesId),
		tgbotapi.NewInlineKeyboardButtonData(TextReject,TextQueryYesId),
	),
)


