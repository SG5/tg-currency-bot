package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
)

type CommandHandler func(u *tgbotapi.Update) string

var (
	chats         = make(map[int64]*Chat)
	commands      = make(map[string]CommandHandler)
	stateCommands = make(map[string]CommandHandler)
)

func startBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_KEY"))

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 600})

	commands["/current"] = currentCurrencyHandler
	commands["текущий курс"] = currentCurrencyHandler

	stateCommands["/min"] = minCurrencyHandler
	stateCommands["минимум"] = minCurrencyHandler
	//stateCommands["/max"] = fromHandler
	//stateCommands["максимум"] = fromHandler

	go func() {
		channel, err := GetUpdateChannel("USD000000TOD")
		if nil != err {
			log.Fatal("can't get moex channel ", err)
		}
		handleMoex(channel, bot)
	}()

	for update := range updates {
		if nil == update.Message || 0 == len(update.Message.Text) {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if _, ok := chats[update.Message.Chat.ID]; !ok {
			chat := Chat{
				ChatId: update.Message.Chat.ID,
			}
			db.Find(&chat, update.Message.Chat.ID)

			chats[update.Message.Chat.ID] = &chat
		}

		go bot.Send(handleMessage(&update))
	}
}

func handleMoex(channel <-chan MOEXRow, bot *tgbotapi.BotAPI ) {
	var chats []Chat
	for row := range channel {
		log.Println("got row ", row)
		db.Where("chat_min <= ?", row.Last).Find(&chats)

		for _, chat := range chats {
			bot.Send(
				tgbotapi.NewMessage(chat.ChatId, fmt.Sprintf("Курс понизился до %f", row.Last)),
			)
		}
	}
}

func handleMessage(u *tgbotapi.Update) tgbotapi.Chattable {
	var handler CommandHandler
	var ok bool

	if handler, ok = commands[u.Message.Text]; !ok {
		if handler, ok = stateCommands[chats[u.Message.Chat.ID].ChatLastCommand]; !ok {
			log.Println("!!", chats[u.Message.Chat.ID].ChatLastCommand)
			handler = helpHandler
		}
	}

	msg := handler(u)
	return tgbotapi.NewMessage(u.Message.Chat.ID, msg)
}
