package telegram

import (
	"log"
	"net/http"
	"net/url"

	"encoding/binary"

	"bytes"
	"encoding/gob"

	"gopkg.in/telegram-bot-api.v4"
)

// TgBot type
type TgBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	users   map[int64]string
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
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("telegram.Init: %s\n", err.Error())
		return err
	}
	t.updates = updates
	return nil
}

func (t *TgBot) getUpdates(todb chan map[string][]byte) {
	var msg *tgbotapi.Message
	for up := range t.updates {
		switch {
		case up.Message != nil:
			msg = up.Message
		case up.ChannelPost != nil:
			msg = up.ChannelPost
		case up.EditedChannelPost != nil:
			msg = up.EditedChannelPost
		case up.EditedMessage != nil:
			msg = up.EditedMessage
		default:
			continue
		}
		if _, ok := t.users[msg.Chat.ID]; !ok {
			id := make([]byte, 8)
			binary.LittleEndian.PutUint64(id, uint64(msg.Chat.ID))
			todb <- map[string][]byte{
				"adduser": []byte(msg.Chat.UserName),
				"chatid":  id,
			}
			t.users[msg.Chat.ID] = msg.Chat.UserName
		}
	}
}

// Run run telegram bot
func (t *TgBot) Run(in, out chan map[string][]byte) {
	out <- map[string][]byte{
		"getusers": []byte("getusers"),
	}
	out <- map[string][]byte{
		"deletemsg": []byte("deletemsg"),
	}
	for message := range in {
		for k, v := range message {
			switch k {
			case "users":
				t.getUsers(v)
			case "sendmsg":
				t.sendMessage(out, v)
			case "deletemsg":
				t.deleteMessages(out, v)
			}
		}
	}
}
func (t *TgBot) getUsers(users []byte) error {
	var res map[int64]string
	buf := bytes.NewBuffer(users)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&res)
	if err != nil {
		log.Printf("telegram.getUsers: %s\n", err.Error())
		return err
	}
	t.users = res
	return nil
}

func (t *TgBot) sendMessage(todb chan map[string][]byte, message []byte) error {
	return nil
}
func (t *TgBot) deleteMessages(todb chan map[string][]byte, messages []byte) error {
	return nil
}
