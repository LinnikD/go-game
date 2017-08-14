package main

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main () {
	bot, err := tgbotapi.NewBotAPI("385402864:AAEmuWbihbSEVV7-8Jy0CDLSeLMcrPpI86s")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}

}

