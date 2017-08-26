package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"time"
	"game"
	"math/rand"
)

func main () {
	bot, _ := tgbotapi.NewBotAPI("385402864:AAEmuWbihbSEVV7-8Jy0CDLSeLMcrPpI86s")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	gameMap := make(map[int64]*game.Game)

	gameChanel := make(chan int64)
	go func() {
		for chatID := range gameChanel {
			go gameMap[chatID].Send("End!")
			delete(gameMap, chatID)
		}
	} ()

	for update := range updates {
		if update.Message == nil {
			continue
		} else {
			chatID := update.Message.Chat.ID

			switch {
			case update.Message.Text == "/start":
				gameMap[chatID] = game.NewGame(bot, chatID)

				gameTime := (5 + rand.Intn(5)) * time.Second

				time.AfterFunc(gameTime, func() {
					gameChanel <- chatID
				})
			default:
				if _, ok := gameMap[chatID]; ok {
					gameMap[chatID].Turn(update)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Pls send /start")
					bot.Send(msg)
				}
			}
		}
	}

}

