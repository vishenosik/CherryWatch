package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot          *tgbotapi.BotAPI
	subscribers  = make(map[int64]bool) // Храним chatID подписчиков
	messageText  = "LIGMA! 🚀"
	pollInterval = 10 * time.Second
)

func CheckSSO() (int, error) {

	url := "http://127.0.0.1:4080/api/ping"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}

func handleChecks() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {

		code, err := CheckSSO()
		if err != nil {
			log.Printf("Ошибка проверки сервиса %v", err)
		}

		for chatID := range subscribers {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("SSO check code: %v", code))
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения для %d: %v", chatID, err)
				// Удаляем отписавшихся или заблокировавших бота
				delete(subscribers, chatID)
			}
		}
	}
}

func main() {
	var err error
	// Замените "YOUR_TELEGRAM_BOT_TOKEN" на ваш реальный токен
	bot, err = tgbotapi.NewBotAPI("7667795695:AAFAiFEQvxm9DPZt-Z3_cUJNZbX_T_6oaf0")
	if err != nil {
		log.Panic(err)
	}

	tgbotapi.NewBotCommandScopeAllPrivateChats()

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Канал для обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// tgbotapi.SetLogger()

	// Запускаем горутину для обработки подписок
	go handleChecks()

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		switch text {
		case "/start", "/subscribe":
			subscribers[chatID] = true
			msg := tgbotapi.NewMessage(chatID, "Вы подписались на рассылку! Каждую минуту вам будет приходить сообщение.")
			bot.Send(msg)
		case "/stop", "/unsubscribe":
			delete(subscribers, chatID)
			msg := tgbotapi.NewMessage(chatID, "Вы отписались от рассылки.")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(chatID, "Используйте /subscribe чтобы подписаться или /unsubscribe чтобы отписаться.")
			bot.Send(msg)
		}
	}
}

// Рассылка сообщений подписчикам
func handleSubscriptions() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Рассылка сообщения %d подписчикам", len(subscribers))
		for chatID := range subscribers {
			msg := tgbotapi.NewMessage(chatID, messageText)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения для %d: %v", chatID, err)
				// Удаляем отписавшихся или заблокировавших бота
				delete(subscribers, chatID)
			}
		}
	}
}
