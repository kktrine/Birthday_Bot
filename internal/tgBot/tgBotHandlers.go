package tgBot

import (
	"birthday_bot/internal/model"
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"strconv"
)

func (b *Bot) handleSignIn(update tgbotapi.Update, username, password string) {
	id := update.Message.Chat.ID
	delete(b.tokens, id)
	fmt.Println(username, password)
	token, userId, err := b.signIn(username, password, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Failed to sign in: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Вход успешно. Запомните ваш id: "+strconv.Itoa(userId)+", он понадобится для добавления информации о себе и подписки на Дни Рождения")
	b.tokens[id] = token
	b.bot.Send(msg)
}

func (b *Bot) handleSignUp(update tgbotapi.Update, username, password string) {
	id := update.Message.Chat.ID
	delete(b.tokens, id)
	err := b.signUp(username, password, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось зарегистрироваться: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Регистрация успешно. Теперь вы можете войти в аккаунт")
	b.bot.Send(msg)
}

func (b *Bot) handleGetEmployees(update tgbotapi.Update) {
	token, ok := b.tokens[update.Message.Chat.ID]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо войти в аккаунт")
		b.bot.Send(msg)
		return
	}
	data, err := b.getEmployees(token)
	id := update.Message.Chat.ID
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось выполнить команду: "+err.Error())
		b.bot.Send(msg)
		return
	}
	var people []model.Employee
	json.Unmarshal(data, &people)
	jsonMessage, err := json.MarshalIndent(people, "", "   ")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(jsonMessage))
	b.bot.Send(msg)
}

func (b *Bot) handleInfo(update tgbotapi.Update, info *model.Employee) {
	token, ok := b.tokens[update.Message.Chat.ID]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо войти в аккаунт")
		b.bot.Send(msg)
		return
	}
	err := b.info(info, token)
	id := update.Message.Chat.ID
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось добавить информацию: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Успешно добавлено")
	b.bot.Send(msg)
}

func (b *Bot) handleSubscribe(update tgbotapi.Update, sub model.Subscribe) {
	token, ok := b.tokens[update.Message.Chat.ID]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо войти в аккаунт")
		b.bot.Send(msg)
		return
	}
	err := b.subscribe(sub, token)
	id := update.Message.Chat.ID
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось добавить информацию: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Подписки успешно добавили")
	b.bot.Send(msg)
}

func (b *Bot) handleUnSubscribe(update tgbotapi.Update, unsubscribe model.Subscribe) {
	token, ok := b.tokens[update.Message.Chat.ID]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо войти в аккаунт")
		b.bot.Send(msg)
		return
	}
	err := b.unsubscribe(unsubscribe, token)
	id := update.Message.Chat.ID
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось добавить информацию: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Подписки успешно добавили")
	b.bot.Send(msg)
}
