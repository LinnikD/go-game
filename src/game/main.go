package game

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"config"
	"model"
	"words"
	"math/rand"
)

func NewGame(bot *tgbotapi.BotAPI, chatID int64, cfg *config.Config) (*Game) {
	patterns := [4]string{"ко", "ро", "ма", "зе"}
	pattern := patterns[rand.Intn(4)]
	mongo := model.NewConnection(cfg.Mongo)

	checker := words.NewWordChecker(cfg.Yandex, mongo)

	g := Game{
		bot: bot,
		chatID: chatID,
		words: make(map[string]struct{}),
		users: make(map[*tgbotapi.User] int),
		pattern: pattern,
		checker: checker,
	}

	g.Send(fmt.Sprintf("You pattern is %s. Go!", g.pattern))
	return &g
}

type Game struct {
	bot      *tgbotapi.BotAPI
	chatID   int64
	words    map[string]struct{}
	users    map[*tgbotapi.User] int
	pattern  string
	checker  *words.WordChecker
}

func (g *Game) Turn(u tgbotapi.Update) {
	message := u.Message.Text
	if len(message) < len(g.pattern) || message[:4] != g.pattern {
		g.Send("YOU LOOOOOSE (not by rules)")
		return
	}
	if _, ok := g.words[message]; ok {
		g.Send("YOU LOOOOOSE (text already was)")
		return
	}
	if g.checker.CheckWordExists(message) {
		g.words[message] = struct{}{}
		if current_points, ok := g.users[u.Message.From]; ok {
			g.users[u.Message.From] = current_points + 1
		} else {
			g.users[u.Message.From] = 1
		}
	} else {
		g.Send("YOU LOOOOOSE (text is not correct)")
	}
}

func (g *Game) Send(text string) {
	msg := tgbotapi.NewMessage(g.chatID, text)
	g.bot.Send(msg)
}

func (g *Game) ShowVictor() {
	max_score := -1
	if len(g.users) == 0 {
		g.Send("End! Wait... You didn't play with me :(((")
		return
	}

	winner := tgbotapi.User{}
	for user, score := range g.users {
		if score > max_score {
			max_score = score
			winner = *user
		}
	}
	winner_name := ""
	if winner.UserName != "" {
		winner_name = winner.UserName
	} else {
		winner_name = winner.FirstName
	}
	g.Send(fmt.Sprintf("End! And the winner is ... @%s!", winner_name))
}
