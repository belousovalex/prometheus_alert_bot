package rocketchat

import (
	"encoding/json"
	"fmt"
	"github.com/belousovalex/prometheus_alert_bot/pkg/core"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Sender struct {
	WebHook string
	filter *core.SameSendFilter
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Attachment struct {
	Title string `json:"title"`
	Fields []Field `json:"fields"`
}

type Message struct {
	Attachments []Attachment `json:"attachments"`
}

func MakeRocketChatSender(webhook string) *Sender{
	return &Sender{
		WebHook:webhook,
		filter:core.MakeSameSendFilter(90 * 60),
	}
}

func (sender *Sender) Send(messages []*core.PrometheusAlert) {
	var fields []Field
	var mss []string
	for _, message := range messages {
		fields = append(fields, Field{
			fmt.Sprintf("%s (%s)", message.Annotations.Summary, message.Status),
			message.Annotations.Description,
		})
		mss = append(
			mss,
			fmt.Sprintf(
				"%s%s%s",
				message.Annotations.Summary,
				message.Status,
				message.Annotations.Description,
			),
		)
	}

	hash := core.MakeHashOfSortingList(mss)
	if sender.filter.IsNeedSend(hash) == false {
		fmt.Printf("Rocketchat. Same message.")
		return
	}
	sender.filter.SetSentTime(hash)


	payload, err := json.Marshal(Message{[]Attachment{{
		Title: "Сообщение об ошибках",
		Fields:fields,
	}}})

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	postData := url.Values{
		"payload": {string(payload)},
	}

	fmt.Printf("Send to rocketchat: %+v\n", postData)
	if resp, err := http.PostForm(sender.WebHook, postData); err != nil {
		fmt.Printf(err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		fmt.Printf("Rocketchat answer: %s\n", string(body))
	}
}
