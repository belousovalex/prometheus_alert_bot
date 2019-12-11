package telegram

import (
	"fmt"
	"github.com/belousovalex/prometheus_alert_bot/pkg/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type TelegramSender struct {
	Token  string
	ChatId string
	filter *core.SameSendFilter
}

func MakeTelegramSender(token string, chatid string) *TelegramSender {
	return &TelegramSender{
		Token:token,
		ChatId:chatid,
		filter:core.MakeSameSendFilter(90 * 60),
	}
}

func (sender *TelegramSender) Send(messages []*core.PrometheusAlert) {
	urlStr := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", sender.Token)
	var messagesStr []string
	for _, message := range messages {
		messageStr := fmt.Sprintf(
			"<b>%s (status=%s)</b>\n%s",
			message.Annotations.Summary,
			message.Status,
			message.Annotations.Description,
		)
		messagesStr = append(messagesStr, messageStr)
	}

	hash := core.MakeHashOfSortingList(messagesStr)
	if sender.filter.IsNeedSend(hash) == false {
		fmt.Printf("Telegram. Same message.")
		return
	}
	sender.filter.SetSentTime(hash)

	postData := url.Values{
		"chat_id":    {sender.ChatId},
		"parse_mode": {"HTML"},
		"text":       {strings.Join(messagesStr, "\n\n")},
	}
	fmt.Printf("Send to telegram: %+v\n", postData)
	if resp, err := http.PostForm(urlStr, postData); err != nil {
		fmt.Printf(err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		fmt.Printf("Telegram answer: %s\n", string(body))
	}
}
