package commands

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/fr0stylo/esveikata-registracija/pkg/api"
	"github.com/fr0stylo/esveikata-registracija/pkg/database"
	"github.com/fr0stylo/esveikata-registracija/pkg/database/entities"
)

func cacheSpecialists(specialistDb *database.Database[*entities.SpecialistItem]) ([]*entities.SpecialistItem, error) {
	prefix := time.Now().Format("2006-01-02")

	items, err := specialistDb.GetByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	if items == nil {
		s, err := api.Get[api.Specialist]("https://ipr.esveikata.lt/api/searchesNew/specialists")
		if err != nil {
			return nil, err
		}

		mapped := lo.Map(s.Data, func(s api.Specialist, id int) *entities.SpecialistItem {
			return &entities.SpecialistItem{ID: s.ID, Name: s.FullName}
		})

		specialistDb.BatchInsert(mapped, func(si *entities.SpecialistItem) ([]byte, error) {
			return []byte(fmt.Sprintf("%s-%d-%s", prefix, si.ID, si.String())), nil
		})

		return mapped, nil
	}

	return items, nil
}

func Find(update tgbotapi.Update, bot *tgbotapi.BotAPI, specialistDb *database.Database[*entities.SpecialistItem]) error {
	mapped, err := cacheSpecialists(specialistDb)
	if err != nil {
		return err
	}

	if update.Message.CommandArguments() != "" {
		mapped = lo.Filter(mapped, func(i *entities.SpecialistItem, _ int) bool {
			return strings.Contains(strings.ToLower(i.Name), strings.ToLower(update.Message.CommandArguments()))
		})
	} else {
		mapped = mapped[:5]
	}

	_, err = bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Radę tinkamą specialistą atsakykite į tą žinutę žodžiu 'Tinka'\n\nSpecialistai:"))
	if err != nil {
		return err
	}
	for _, v := range mapped {
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.From.ID, v.String()))
		if err != nil {
			return err
		}
	}

	return nil
}
