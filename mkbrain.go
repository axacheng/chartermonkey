package main

import (
	"chartermonkey/mknote"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {

	mknote.InitDB()

	mknote.Query()

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req) //events is []*Event
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events { //event is *Event
			if event.Type == linebot.EventTypeMessage { //"message" should be defined in /callback
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "恰特猴" {
						message.Text = "幹嘛~?"
					} else if message.Text == "list" {
						message.Text = mknote.Query()
					} else if message.Text == "+1" && event.Source.GroupID != "" {
						profile, err := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do()
						if err != nil {
							log.Print(err)
						}
						message.Text = "好喔, " + profile.DisplayName + " +1, 吱吱"
					} else if message.Text == "+1" && event.Source.GroupID == "" {
						profile, err := bot.GetProfile(event.Source.UserID).Do()
						if err != nil {
							log.Print(err)
						}
						date := event.Postback.Params.Date
						message.Text = "好喔, 今天是" + date + ", 下次 " + profile.DisplayName + " +1, 吱吱"
					}
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do()
					if err != nil {
						log.Print(err)
					}
				}

			}
		}
	})
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
