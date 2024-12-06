package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"progect/controllers"
)

func startTelegramBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil {
				// логируем Chat ID и текст сообщения
				log.Printf("Chat ID: %d, Username: [%s], Message: %s",
					update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)

				// текст для ответа пользователю
				var replyText string
				switch update.Message.Text {
				case "/start":
					replyText = fmt.Sprintf("Добро пожаловать! Ваш уникальный Chat ID: %d", update.Message.Chat.ID)
				case "/help":
					replyText = "Список команд: /start, /help, ..."
				default:
					replyText = fmt.Sprintf("Привет! Вы написали: %s\nВаш Chat ID: %d",
						update.Message.Text, update.Message.Chat.ID)
				}
				// отправка сообщения пользователю
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
				if _, err := bot.Send(reply); err != nil {
					log.Printf("Ошибка отправки сообщения: %v", err)
				}
			}
		}
	}()

}

func main() {
	// запуск Telegram-бота
	botToken := "7693218177:AAEqxnmF1W4OvegCx0g6mWAwB_EpwppXPnw"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка при инициализации Telegram API: %v", err)
	}
	log.Printf("Запущен бот %s", bot.Self.UserName)

	controllers.TelegramBot = bot

	// запуск обработчика Telegram-бота
	startTelegramBot(bot)

	// настройка маршрутов веб-сервера
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static"))))
	http.HandleFunc("/", controllers.FormHandler)
	http.HandleFunc("/Portfolio", controllers.FormHandlerP)

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
