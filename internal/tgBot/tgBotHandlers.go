package tgBot

import (
	"birthday_bot/internal/model"
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func (b *Bot) handleSignIn(update tgbotapi.Update) {
	id := update.Message.Chat.ID
	if b.states[id] != "ok" {
		return
	}
	username, _ := b.usernames[id]
	password, _ := b.passwords[id]
	fmt.Println(username, password)
	token, err := b.signIn(username, password)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to sign in: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Successfully signed in. Your token: "+token+"\nToken will be sent automatically, you don't need to keep it")
	b.tokens[update.Message.Chat.ID] = token
	b.bot.Send(msg)
}

func (b *Bot) handleSignUp(update tgbotapi.Update) {
	id := update.Message.Chat.ID
	username, _ := b.usernames[id]
	password, _ := b.passwords[id]
	err := b.signUp(username, password)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to sign up: "+err.Error())
		b.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Successfully signed up. You can now sign in.")
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
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось выполнить команду: "+err.Error())
		b.bot.Send(msg)
		return
	}
	var people []model.Employee
	json.Unmarshal(data, &people)
	jsonMessage, err := json.MarshalIndent(people, "", "   ")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(jsonMessage))
	b.bot.Send(msg)
}
