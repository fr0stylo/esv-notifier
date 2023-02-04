package commands

import (
	"bytes"
	"encoding/gob"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/fr0stylo/esveikata-registracija/pkg/database"
	"github.com/fr0stylo/esveikata-registracija/pkg/database/entities"
)

func stopAll(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *database.Database[*entities.Registrations], updateKeyMap [][]byte) error {
	if err := db.BatchUpdate(updateKeyMap, func(value *entities.Registrations) ([]byte, error) {
		filtered := lo.Filter(*value, func(v int64, _ int) bool {
			return v != update.Message.From.ID
		})

		buf := bytes.NewBuffer(nil)
		if err := gob.NewEncoder(buf).Encode(filtered); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}); err != nil {
		return err
	}

	return nil
}

func Stop(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *database.Database[*entities.Registrations]) error {
	updateKeyMap := make([][]byte, 0)
	if err := db.Iterate(func(key string, value *entities.Registrations) error {
		if lo.Contains(*value, update.Message.From.ID) {
			updateKeyMap = append(updateKeyMap, []byte(key))
		}
		return nil
	}); err != nil {
		return err
	}

	if update.Message.CommandArguments() == "visus" {
		return stopAll(update, bot, db, updateKeyMap)
	}

	return nil
}
