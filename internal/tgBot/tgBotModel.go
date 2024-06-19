package tgBot

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
)

type Bot struct {
	apiBaseURL string
	bot        *tgbotapi.BotAPI
	updates    tgbotapi.UpdatesChannel
	tokens     map[int64]string
	states     map[int64]string
	usernames  map[int64]string
	passwords  map[int64]string
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Token string `json:"token"`
}

func NewBot(apiBaseURL string) *Bot {
	bot, err := tgbotapi.NewBotAPI("6264760242:AAGga9UN4U4cditVpvKTc7mWRNo5nAbyNP4")
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	return &Bot{apiBaseURL: apiBaseURL, bot: bot, updates: updates, tokens: make(map[int64]string), states: make(map[int64]string)}
}

func (b *Bot) Start() {
	flag := true
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Please use /signin or /signup to get started.")
			b.bot.Send(msg)
		case "signin":
			b.handleSignIn(update)
		case "signup":
			b.handleSignUp(update)
		case "employees":
			b.handleGetEmployees(update)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
			b.bot.Send(msg)
		}
	}
}
