package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"time"
	"game"
	"math/rand"
	"config"
	"log"
)

func main () {
	cfg, err := config.LoadConfig("/home/uzzz/go/go-game/src/config/config.yaml")
	if err != nil {
		log.Panicln("LoadConfig: ", err)
	}

	bot, _ := tgbotapi.NewBotAPI(cfg.Token)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	gameMap := make(map[int64]*game.Game)

	gameChanel := make(chan int64)
	go func() {
		for chatID := range gameChanel {
			go gameMap[chatID].ShowVictor()
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
				gameMap[chatID] = game.NewGame(bot, chatID, &cfg)

				gameTime := time.Duration(rand.Intn(20 - 5) + 5) * time.Second
				gameTime = 10 * time.Second
				time.AfterFunc(gameTime, func() {
					gameChanel <- chatID
				})
			default:
				if _, ok := gameMap[chatID]; ok {
					gameMap[chatID].Turn(update)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Хочешь знать кто крче батлит? Жми /start")
					bot.Send(msg)
				}
			}
		}
	}

}

