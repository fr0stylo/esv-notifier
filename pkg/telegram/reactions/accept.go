package reactions

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/fr0stylo/esveikata-registracija/pkg/database"
	"github.com/fr0stylo/esveikata-registracija/pkg/database/entities"
)

func Accept(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *database.Database[*entities.Registrations]) error {
	if update.Message.Text == "Tinka" && update.Message.ReplyToMessage != nil && []rune(update.Message.ReplyToMessage.Text)[0] == '(' {
		id := strings.Split(update.Message.ReplyToMessage.Text, " ")[0]
		id = strings.Replace(id, "(", "", -1)
		id = strings.Replace(id, ")", "", -1)

		if err := db.Upsert(id, update.Message.From.ID); err != nil {
			return err
		}
	}

	return nil
}
