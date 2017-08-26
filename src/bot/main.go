package main

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"time"
)

func newGame(bot *tgbotapi.BotAPI, chatID int64) (*Game) {
	game := Game{bot, chatID, make(map[string]struct{})}
	game.Send("Started!")
	return &game
}

type Game struct {
	bot *tgbotapi.BotAPI
	chatID int64
	words map[string]struct{}
}

func (g *Game) Turn(u tgbotapi.Update) {
	message := u.Message.Text

	if _, ok := g.words[message]; ok {
		g.Send("YOU LOOOOOSE")
	} else {
		g.words[message] = struct {}{}
	}
}

func (g *Game) Send(text string) {
	msg := tgbotapi.NewMessage(g.chatID, text)
	g.bot.Send(msg)
}

func main () {
	bot, _ := tgbotapi.NewBotAPI("385402864:AAEmuWbihbSEVV7-8Jy0CDLSeLMcrPpI86s")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1
	updates, _ := bot.GetUpdatesChan(u)

	gameMap := make(map[int64]*Game)

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
				gameMap[chatID] = newGame(bot, chatID)

				time.AfterFunc(15 * time.Second, func() {
					gameChanel <- chatID
				})
			default:
				if _, ok := gameMap[chatID]; ok {
					gameMap[chatID].Turn(update)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Pls start o/")
					bot.Send(msg)
				}
			}

			//msg := tgbotapi.NewMessage(ChatID, "Heil o/")
			//bot.Send(msg)


		}
	}

}

