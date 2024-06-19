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
	Username string `json:"username"`
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
	return &Bot{apiBaseURL: apiBaseURL, bot: bot, updates: updates, tokens: make(map[int64]string), states: make(map[int64]string), usernames: make(map[int64]string), passwords: make(map[int64]string)}
}

func (b *Bot) Start() {
	flag := false
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		var state string
		id := update.Message.Chat.ID
		state, ok := b.states[id]
		if !ok {
			state = "needUsername"
			b.states[id] = state
		}
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(id, "Welcome! Please use /signin or /signup to get started.")
			b.bot.Send(msg)
			flag = false
		case "signin":
			if state != "ok" {
				b.waitForLoginPassword(id)
				flag = true
			} else {
				b.handleSignIn(update)
				flag = false
			}
		case "signup":
			if state != "ok" {
				b.waitForLoginPassword(id)
				flag = true
			} else {
				b.handleSignUp(update)
				flag = false
			}
		case "employees":
			b.handleGetEmployees(update)
			flag = false
		default:
			if !flag || state == "ok" {
				msg := tgbotapi.NewMessage(id, "Неизвестная команда")
				flag = false
				b.bot.Send(msg)
			} else if state == "needUsername" {
				b.usernames[id] = update.Message.Text
				b.states[id] = "needPassword"
				msg := tgbotapi.NewMessage(id, "Введите пароль")
				b.bot.Send(msg)
				flag = true
			} else {
				b.passwords[id] = update.Message.Text
				b.states[id] = "ok"
				msg := tgbotapi.NewMessage(id, "Выберите команду")
				b.bot.Send(msg)
				flag = false
			}
		}
	}
}

func (b *Bot) waitForLoginPassword(id int64) {
	switch b.states[id] {
	case "needUsername":
		msg := tgbotapi.NewMessage(id, "Введите логин")
		b.bot.Send(msg)
	case "needPassword":
		msg := tgbotapi.NewMessage(id, "Введите пароль")
		b.bot.Send(msg)
	}
}
