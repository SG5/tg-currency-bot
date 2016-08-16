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
	if "/from" == u.Message.Text {
		text = "Откуда выезжаешь?"
	} else if "/to" == u.Message.Text {
		text = "Куда хочешь поехать?"
	} else if "/date" == u.Message.Text {
		text = "Когда хочешь выехать? Используй формат 31.01.2016"
	} else if "/price" == u.Message.Text {
		text = "Сколько есть денег на билет?"
	} else if "/help" == u.Message.Text {
		text = "/from - Город отправления\n" +
			"/to - Город прибытия\n" +
			"/date - Дата отправления\n" +
			"/price - Максимальная цена билета\n" +
			"/next - Найти ближайший поезд\n" +
			"/info - Показать твои настроики\n" +
			"/cancel - Отмена последней команды\n" +
			"/help - Список доступных команд\n"
		chats[u.Message.Chat.ID].CancelCommand()
	} else {
		text = "Набери /help для списка команд"
	}
	return text
}
