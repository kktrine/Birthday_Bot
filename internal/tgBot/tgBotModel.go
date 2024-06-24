package tgBot

import (
	"birthday_bot/internal/model"
	"birthday_bot/internal/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
	"time"
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

//type signUpResponce struct {
//	Token string `json:"token"`
//	Id    int    `json:"id"`
//}

type signInResponse struct {
	Id    int    `json:"id"`
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

func (b *Bot) Start(db *storage.Storage) {
	go b.periodSend(db)
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		id := update.Message.Chat.ID
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(id, "Привет! Выбери команду в меню.")
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
		case "info":
			info := b.getInfo(id)
			if info != nil {
				b.handleInfo(update, info)
			} else {
				msg := tgbotapi.NewMessage(id, "Ошибка: данные не заполнены")
				b.bot.Send(msg)
			}
		case "subscribe":
			ids := b.getIds(id)
			if ids != nil {
				b.handleSubscribe(update, *ids)
			} else {
				msg := tgbotapi.NewMessage(id, "Ошибка: данные не заполнены")
				b.bot.Send(msg)
			}
		case "unsubscribe":
			ids := b.getIds(id)
			if ids != nil {
				b.handleUnSubscribe(update, *ids)
			} else {
				msg := tgbotapi.NewMessage(id, "Ошибка: данные не заполнены")
				b.bot.Send(msg)
			}
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

func (b *Bot) getInfo(id int64) *model.Employee {
	res := model.Employee{}
	msg := tgbotapi.NewMessage(id, "Введите id или /exit, чтобы выйти")
	b.bot.Send(msg)
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Command() == "exit" {
			return nil
		} else if update.Message.Command() != "" {
			msg := tgbotapi.NewMessage(id, "Неверная команда.\nВведите /exit, если хотите выйти")
			b.bot.Send(msg)
		}
		if res.UserId <= 0 {
			idString := update.Message.Text
			idInt, err := strconv.Atoi(idString)
			if err != nil {
				msg := tgbotapi.NewMessage(id, "Id не распознан. Введите id.\nВведите /exit, если хотите выйти")
				b.bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(id, "Id принят. Введите имя.\nВведите /exit, если хотите выйти")
				b.bot.Send(msg)
				res.UserId = idInt
			}
		} else if res.Name == "" {
			res.Name = update.Message.Text
			msg := tgbotapi.NewMessage(id, "Имя принято.\nВведите фамилию или /exit, чтобы выйти")
			b.bot.Send(msg)
		} else if res.Surname == "" {
			res.Surname = update.Message.Text
			msg := tgbotapi.NewMessage(id, "Фамилия принята.\nВведите дату рождения в формате дд.мм.гггг или /exit, чтобы выйти")
			b.bot.Send(msg)
		} else if res.Birth == nil {
			dateStr := update.Message.Text
			parsedDate, err := time.Parse("02.01.2006", dateStr)
			if err != nil {
				msg := tgbotapi.NewMessage(id, "Неверный ввод.\nВведите дату рождения в формате дд.мм.гггг или /exit, чтобы выйти")
				b.bot.Send(msg)

			} else {
				msg := tgbotapi.NewMessage(id, "Дата рождения принята")
				b.bot.Send(msg)
				res.Birth = &parsedDate
				return &res

			}

		}
	}
	return &res
}

func (b *Bot) getIds(id int64) *model.Subscribe {
	res := model.Subscribe{}
	ids := make([]int, 0)
	msg := tgbotapi.NewMessage(id, "Введите Ваш id или /exit, чтобы закончить ввод")
	b.bot.Send(msg)
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Command() == "exit" {
			if len(ids) > 0 {
				res.SubscribeTo = &ids
				return &res
			}
			return nil
		} else if update.Message.Command() != "" {
			msg := tgbotapi.NewMessage(id, "Неверная команда.\nВведите /exit, если хотите закончить ввод")
			b.bot.Send(msg)
		}
		currentIdString := update.Message.Text
		currentIdInt, err := strconv.Atoi(currentIdString)
		if err != nil {
			msg := tgbotapi.NewMessage(id, "Неверный id.\nВведите верный id или /exit, если хотите закончить ввод")
			b.bot.Send(msg)
		} else if res.Id == nil {
			res.Id = &currentIdInt
			msg := tgbotapi.NewMessage(id, "id принят.\nВведите id пользователя или /exit, чтобы выйти")
			b.bot.Send(msg)
		} else {
			ids = append(ids, currentIdInt)
			msg := tgbotapi.NewMessage(id, "id принят.\nВведите еще id или /exit, чтобы выйти")
			b.bot.Send(msg)
		}
	}
	return &res
}

func (b *Bot) periodSend(db *storage.Storage) {
	for {
		now := time.Now()
		//next := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location()).Add(24 * time.Hour)
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()).Add(1 * time.Minute)
		pause := next.Sub(now)
		time.Sleep(pause)
		b.sendNotifications(db)
	}
}

func (b *Bot) sendNotifications(db *storage.Storage) {
	messages := db.CheckBDays()
	if messages == nil {
		return
	}
	for _, message := range *messages {
		for _, person := range message.BdayPeople {
			msg := tgbotapi.NewMessage(message.ChatId, "Today is "+person+"'s Birthday!")
			b.bot.Send(msg)
		}
	}
}
