package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	for ; ; time.Sleep(time.Minute * 1) {
		log.Println("checking...")
		checking()
		log.Println("done checking")
	}
}
func init() {
	http.DefaultClient.Timeout = time.Second * 3
}

func checking() {
	var target = os.Getenv("target")
	resp, err := http.DefaultClient.Get(target)
	if err != nil {
		alert(fmt.Sprintf("health check to %s fail with error: %s", target, err.Error()))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		alert(fmt.Sprintf("health check to %s fail with status: %s", target, resp.Status))
	}
}
func alert(msg string) {
	telegramIDStr := os.Getenv("telegram_server_alert_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		log.Println("parse telegram id fail: " + err.Error())
		return
	}
	if err := sendTelegram(telegramID, msg); err != nil {
		log.Println("health check fail: " + err.Error())
	}
}

func sendTelegram(id int64, msg string) error {
	log.Printf("alerting: %s\n", msg)
	telegramToken := os.Getenv("telegram_server_alert_token")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}
	msgReq := tgbotapi.NewMessage(id, msg)
	_, err = bot.Send(msgReq)
	if err != nil {
		log.Printf("alerting to telegram fail: %v\n", err)
		return err
	}
	return nil
}
