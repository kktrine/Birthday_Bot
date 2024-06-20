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
	token, err := b.signIn(username, password, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Failed to sign in: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Successfully signed in. Your token: "+token+"\nToken will be sent automatically, you don't need to keep it")
	b.tokens[id] = token
	b.bot.Send(msg)
}

func (b *Bot) handleSignUp(update tgbotapi.Update, username, password string) {
	id := update.Message.Chat.ID
	delete(b.tokens, id)
	userId, err := b.signUp(username, password, id)
	if err != nil {
		msg := tgbotapi.NewMessage(id, "Не удалось зарегистрироваться: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(id, "Регистрация успешно. Ваш id: "+strconv.Itoa(userId)+"\nТеперь вы можете войти в аккаунт")
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
