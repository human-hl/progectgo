package controllers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"html/template"
	"log"
	"net/http"
	"progect/models"
)

var TelegramBot *tgbotapi.BotAPI

var Photos = [...]string{"blackfirst", "blacksecond", "colorthird", "Polinafirst"}

func InitTelegramBot(token string) {
	var err error
	TelegramBot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка при инициализации Telegram-бота: %v", err)
	}
	log.Printf("Telegram-бот успешно запущен: %s", TelegramBot.Self.UserName)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	dsn := "root:12345@tcp(localhost:3306)/Send_txtemail"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка подключения к базе данных: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при проверке подключения к базе данных: %v", err), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		fio := r.FormValue("FIO")
		email := r.FormValue("email")
		price := r.FormValue("price")
		if fio == "" || email == "" || price == "" {
			http.Error(w, "Все поля должны быть заполнены", http.StatusBadRequest)
			return
		}

		texte, err := models.GetTextEmail(db)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Ошибка получения текста для email", http.StatusInternalServerError)
			return
		}

		var msg []byte
		if price == "монтаж видео" {
			for _, text := range texte {
				if text.ID == 2 {
					msg = []byte(fmt.Sprintf("%s %s %s", fio, text.DESCRIPTION, email))
					break
				}
			}
		} else if price == "ретуш фото" {
			for _, text := range texte {
				if text.ID == 3 {
					msg = []byte(fmt.Sprintf("%s %s %s", fio, text.DESCRIPTION, email))
					break
				}
			}
		} else {
			for _, text := range texte {
				if text.ID == 1 {
					msg = []byte(fmt.Sprintf("%s %s %s %s", fio, text.DESCRIPTION, price, email))
					break
				}
			}
		}

		if msg == nil {
			http.Error(w, "Ошибка: сообщение не сформировано", http.StatusInternalServerError)
			return
		}

		// отправка сообщения в Telegram
		if TelegramBot != nil {
			tgMessage := tgbotapi.NewMessage(, string(msg))
			_, err := TelegramBot.Send(tgMessage)
			if err != nil {
				log.Printf("Ошибка отправки сообщения в Telegram: %v", err)
			} else {
				log.Println("Сообщение успешно отправлено в Telegram")
			}
		} else {
			log.Println("Telegram-бот не инициализирован")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./templates/layout.html", "./templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обработке шаблона: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "My Go Server",
	}
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при выполнении шаблона: %v", err), http.StatusInternalServerError)
	}
}

// <div class="column">
//
//	    <img src="/static/images/blacksecond.svg">
//	    <img src="/static/images/blackfirst.svg">
//	    <img src="/static/images/Polinafirst.svg">
//	    <img src="/static/images/colorthird.svg">
//	</div>
// func PostPhoto(photo []string, count int) {
// 	photos := `<div class="column">
// 	<img src="/static/images/%s">
// 	<img src="/static/images/%s">
// 	<img src="/static/images/%s">
// 	<img src="/static/images/%s">
// 	</div>`
// 	var photostwo []byte
// 	for i := 0; i < count; i++ {
// 		for i := 0; i < len(photo); i++ {
// 			photostwo = []byte(fmt.Sprintf(photos, i, i+1, i+2, i+3))
// 		}

// 	}
// }

func FormHandlerP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/layout.html", "./templates/portfolio.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обработке шаблона: %v", err), http.StatusInternalServerError)
		return
	}
	data := struct {
		Title string
	}{
		Title: "My Go Server",
	}
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при выполнении шаблона: %v", err), http.StatusInternalServerError)
	}
}
