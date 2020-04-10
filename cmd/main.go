package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	"flag"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"net/http"
	"os"
)


var (
	hostF = flag.String("host", "", "Server host")
	portF = flag.String("port", "", "Server port")
)

func main () {
	log.Print("start application")
	flag.Parse()

	host, ok := FlagOrEnv(*hostF, betypes.EnvHost)
	if !ok {
		log.Panic("can't get host")
	}

	port, ok := FlagOrEnv(*portF, betypes.EnvPort)
	if !ok {
		log.Panic("can't get port")
	}

	addr := net.JoinHostPort(host, port)
	log.Println(host,port)

	start(addr)
}




func start (addr string) {

	db := dataBase.Connect()
	bot := botConnect()

	go log.Fatal(http.ListenAndServe(addr,  nil))

	admin := betypes.User{
		ChatId: betypes.AdminChatId,
		FirstName: betypes.AdminFirstName,
		Phone:betypes.AdminPhone,
		Role:betypes.AdminRole,
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	for update := range updates {
		log.Println(dataBase.IsUserExist(db,admin))
		log.Printf("%+v\n", update)
	}
}


func botConnect () *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(betypes.BotToken)
	if err != nil {
		log.Fatal("Can't connect to bot api")
	}
	log.Printf("Authorized bot api - %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(betypes.BotWebhook))
	if err != nil {
		log.Fatal("Can't connect set webhook of telegram-bot")
	}
	return bot
}


func FlagOrEnv(flag string, envKey string) (string, bool) {
	if flag != "" {
		return flag, true
	}
	return os.LookupEnv(envKey)
}