package telegram

import (
	"log"
	"net/http"
	"net/url"

	"gopkg.in/telegram-bot-api.v4"
)

// TgBot type
type TgBot struct {
	bot *tgbotapi.BotAPI
}

// Init initialise TgBot
func (t *TgBot) Init(conf *map[string]string) error {
	var token string
	var httpclient *http.Client
	var withproxy bool
	for k, v := range *conf {
		switch k {
		case "token":
			token = v
		case "tgproxy":
			proxy, _ := url.Parse(v)
			httpclient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
			withproxy = true
		}
	}
	if withproxy {
		bot, err := tgbotapi.NewBotAPIWithClient(token, httpclient)
		if err != nil {
			log.Printf("telegram.Init: %s\n", err.Error())
			return err
		}
		t.bot = bot
	} else {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Printf("telegram.Init: %s\n", err.Error())
			return err
		}
		t.bot = bot
	}
}
