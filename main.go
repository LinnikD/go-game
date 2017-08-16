package main

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"time"
)

func main () {
	bot, _ := tgbotapi.NewBotAPI("385402864:AAEmuWbihbSEVV7-8Jy0CDLSeLMcrPpI86s")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Heil")
			bot.Send(msg)

			time.AfterFunc(3 * time.Second, func() {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Bye")
				bot.Send(msg)
			})
		}
	}

}

