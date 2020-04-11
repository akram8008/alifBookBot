package main

import (
	"alifLibrary/betypes"
	dataBase "alifLibrary/crud"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)


/*var (
	hostF = flag.String("host", "", "Server host")
	portF = flag.String("port", "", "Server port")
)
*/
func main () {
	log.Print("start application")
	//flag.Parse()

	/*host, ok := FlagOrEnv(*hostF, betypes.EnvHost)
	if !ok {
		log.Panic("can't get host")
	}
		addr := net.JoinHostPort(host, port)
	log.Println(host,port)

	*/

	port := os.Getenv(betypes.EnvPort)

	start(":"+port)
}




func start (addr string) {
	log.Println(addr)
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	log.Println("OK")
	db := dataBase.Connect()
	bot := botConnect(addr)



	updates := bot.ListenForWebhook("/")
	for update := range updates {
	if _,err :=bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,update.Message.Text));err!=nil {
		admin := betypes.User{ChatId:461795511}
		admin,ok,err := dataBase.IsUserExist(db,admin)
		log.Println(admin,ok,err)
	}
		/*
		if update.Message != nil {
			newMessage(update,db)
		}
		if update.CallbackQuery != nil {
			//newQuery(update,db)
		}
	*/
	}
}



func botConnect (addr string) *tgbotapi.BotAPI {
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


func FlagOrEnv(flag string, envKey string) (string, bool) {
	if flag != "" {
		return flag, true
	}
	return os.LookupEnv(envKey)
}