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
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ChatId   int64  `json:"chatId"`
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
	return &Bot{apiBaseURL: apiBaseURL, bot: bot, updates: updates, tokens: make(map[int64]string)}
}

func (b *Bot) Start() {
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		id := update.Message.Chat.ID
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(id, "Привет! Выбери команду:\n/signin\n/signup\n/employees")
			b.bot.Send(msg)
		case "signin":
			user, password := b.getUsernamePassword(id)
			if user != "" && password != "" {
				b.handleSignIn(update, user, password)
			} else {
				msg := tgbotapi.NewMessage(id, "Ошибка входа: данные не заполнены")
				b.bot.Send(msg)
			}
		case "signup":
			user, password := b.getUsernamePassword(id)
			if user != "" && password != "" {
				b.handleSignUp(update, user, password)
			} else {
				msg := tgbotapi.NewMessage(id, "Ошибка регистрации: данные не заполнены")
				b.bot.Send(msg)
			}
		case "employees":
			b.handleGetEmployees(update)
		default:
			msg := tgbotapi.NewMessage(id, "Неизвестная команда")
			b.bot.Send(msg)
		}
	}
}

func (b *Bot) getUsernamePassword(id int64) (string, string) {
	username := ""
	password := ""
	msg := tgbotapi.NewMessage(id, "Введите имя пользователя или /exit, чтобы выйти")
	b.bot.Send(msg)
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Command() == "exit" {
			return username, password
		} else if update.Message.Command() != "" {
			msg := tgbotapi.NewMessage(id, "Неверная команда.\nВведите /exit, если хотите выйти")
			b.bot.Send(msg)
		}
		if username == "" {
			username = update.Message.Text
			msg := tgbotapi.NewMessage(id, "Имя пользователя принято.\nВведите пароль или /exit, чтобы выйти")
			b.bot.Send(msg)
		} else {
			password = update.Message.Text
			msg := tgbotapi.NewMessage(id, "Пароль принят")
			b.bot.Send(msg)
			return username, password
		}
	}
	return username, password
}
