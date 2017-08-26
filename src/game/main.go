package game

import (
	"gopkg.in/telegram-bot-api.v4"
)

func NewGame(bot *tgbotapi.BotAPI, chatID int64) (*Game) {
	g := Game{bot, chatID, make(map[string]struct{}), make(map[*tgbotapi.User] int)}
	g.Send("Started!")
	return &g
}

type Game struct {
	bot *tgbotapi.BotAPI
	chatID int64
	words map[string]struct{}
	users map[*tgbotapi.User] int
}

func (g *Game) Turn(u tgbotapi.Update) {
	message := u.Message.Text
	if _, ok := g.words[message]; ok {
		g.Send("YOU LOOOOOSE")
	} else {
		g.words[message] = struct {}{}
		if current_points, ok := g.users[u.Message.From]; ok {
			g.users[u.Message.From] = current_points + 1
		} else {
			g.users[u.Message.From] = 1
		}
	}
}

func (g *Game) Send(text string) {
	msg := tgbotapi.NewMessage(g.chatID, text)
	g.bot.Send(msg)
}
