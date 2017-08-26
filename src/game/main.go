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
	patterns := [11]string{"ко", "ро", "ма", "зе", "ре", "ба", "ла", "да", "ди", "ли", "пр"}
	pattern := patterns[rand.Intn(11)]
	mongo := model.NewConnection(cfg.Mongo)

	checker := words.NewWordChecker(cfg.Yandex, mongo)

	g := Game{
		bot: bot,
		chatID: chatID,
		words: make(map[string]struct{}),
		users_score: make(map[int] int),
		users_names: make(map[int] string),
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
	users_score    map[int] int
	users_names    map[int] string
	pattern  string
	checker  *words.WordChecker
}

func (g *Game) Turn(u tgbotapi.Update) {
	message := u.Message.Text
	user_name := ""
	if u.Message.From.UserName != "" {
		user_name = u.Message.From.UserName
	} else {
		user_name = u.Message.From.FirstName
	}
	if len(message) < len(g.pattern) || message[:4] != g.pattern {
		g.Send(fmt.Sprintf("@%s LOOOOOSE (not by rules)", user_name))
		return
	}
	if _, ok := g.words[message]; ok {
		g.Send(fmt.Sprintf("@%s LOOOOOSE (text already was)", user_name))
		return
	}
	if g.checker.CheckWordExists(message) {
		g.words[message] = struct{}{}
		if current_points, ok := g.users_score[u.Message.From.ID]; ok {
			g.users_score[u.Message.From.ID] = current_points + 1
		} else {
			g.users_score[u.Message.From.ID] = 1
			if u.Message.From.UserName != "" {
				g.users_names[u.Message.From.ID] = u.Message.From.UserName
			} else {
				g.users_names[u.Message.From.ID] = u.Message.From.FirstName
			}
		}
	} else {
		g.Send(fmt.Sprintf("@%s LOOOOOSE (text is not correct)", user_name))
	}
}

func (g *Game) Send(text string) {
	msg := tgbotapi.NewMessage(g.chatID, text)
	g.bot.Send(msg)
}

func (g *Game) ShowVictor() {
	max_score := -1
	if len(g.users_score) == 0 {
		g.Send("End! Wait... You didn't play with me :(((")
		return
	}

	winner := ""
	tableOfResult := "Scores table:\n"
	for user_id, score := range g.users_score {
		user_name := g.users_names[user_id]
		tableOfResult = tableOfResult + fmt.Sprintf("%s: %d\n", user_name, score)
		if score > max_score {
			max_score = score
			winner = user_name
		}

	}
	g.Send(fmt.Sprintf("End! And the winner is ... @%s!\n%s", winner, tableOfResult))
}
