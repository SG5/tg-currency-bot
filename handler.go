package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
)

func currentCurrencyHandler(u *tgbotapi.Update) string {
	return fmt.Sprintf("Текущий биржевой курс доллара: %f рублей", getLastRow().Last)
}

func minCurrencyHandler(u *tgbotapi.Update) string {
	price, err := strconv.ParseFloat(u.Message.Text, 32)
	if nil != err {
		log.Print(err)
		return "Ты уверен что отправил число?"
	}
	chats[u.Message.Chat.ID].ChatMin = float32(price)
	chats[u.Message.Chat.ID].CancelCommand()

	return fmt.Sprintf("Я сообщу тебе когда доллар опустится ниже %f", price)
}

func helpHandler(u *tgbotapi.Update) string {
	if '/' == u.Message.Text[0] {
		chats[u.Message.Chat.ID].ChatLastCommand = u.Message.Text
		chats[u.Message.Chat.ID].Save()
	}

	var text string
	if "/min" == u.Message.Text {
		text = "Укажите цену, к примеру 65.32, и я напишу тебе когда курс опустится ниже"
	} else {
		text = "Набери /help для списка команд"
	}
	return text
}
