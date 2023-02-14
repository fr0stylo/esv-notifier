package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/fr0stylo/esveikata-registracija/pkg/database"
	"github.com/fr0stylo/esveikata-registracija/pkg/database/entities"
	"github.com/fr0stylo/esveikata-registracija/pkg/telegram/commands"
	"github.com/fr0stylo/esveikata-registracija/pkg/telegram/reactions"
)

const ID = 2121299350

func main() {
	db, err := database.NewDatabases("esveikata-registracija.db")
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("token"))
	if err != nil {
		log.Fatal(err)
	}
	defer bot.Send(tgbotapi.NewMessage(ID, "Closing"))

	bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Started %s", time.Now().Format("2006-01-02 15:04:05"))))

	go commandReceiver(bot, db)
	go ticker(bot, db)

	sigchnl := make(chan os.Signal, 1)

	signal.Notify(sigchnl, syscall.SIGTERM, syscall.SIGINT)
	<-sigchnl
}

func ticker(bot *tgbotapi.BotAPI, db *database.Databases) {
	// ticker := time.NewTicker(15 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		db.Registrations().Iterate(func(key string, value *entities.Registrations) error {
			fmt.Printf("key: %s, value: %v\n", key, value)
			// resp, err := api.Get[api.Specialist](fmt.Sprintf("https://ipr.esveikata.lt/api/searches/appointments/times?pecialistId=%s&page=0&size=50", key))
			// if err != nil {
			// 	bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Error: %s", err)))
			// 	return err
			// }

			// if len(resp.Data) > 0 {
			// 	for _, id := range *value {
			// 		bot.Send(tgbotapi.NewMessage(id, fmt.Sprintf("Laisvų laikų:  %d", len(resp.Data))))
			// 	}
			// } else {
			// 	log.Print("Found no results")
			// }
			return nil
		})
	}
}

func commandReceiver(bot *tgbotapi.BotAPI, db *database.Databases) {
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.ReplyToMessage != nil {
			switch update.Message.Text {
			case "Tinka":
				handleError(bot, reactions.Accept(update, bot, db.Registrations()))

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Nezinoma komanda"))
			}

			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "ieskoti":
			handleError(bot, commands.Find(update, bot, db.Specialists()))
		case "stabdyti":
			handleError(bot, commands.Stop(update, bot, db.Registrations()))

			// TODO: add stop command
			// [ ] list all specialists that is currently running
			// [ ] select specialist to stop
			// [*] select all specialists to stop
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Nezinoma komanda"))
		}
	}
}

func handleError(bot *tgbotapi.BotAPI, err error) {
	if err != nil {
		fmt.Println(err)
		bot.Send(tgbotapi.NewMessage(ID, err.Error()))
	}
}
