package services

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegramMessage(
	chatID int64,
	message string,
) error {

	botToken := os.Getenv("BOT_TOKEN")

	fmt.Println("BOT TOKEN EXISTS:", botToken != "")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, message)

	result, err := bot.Send(msg)
	if err != nil {
		return err
	}

	fmt.Println("telegram sent:", result.MessageID)

	return nil
}