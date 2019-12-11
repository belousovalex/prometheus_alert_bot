package sender_factory

import (
	"fmt"
	"github.com/belousovalex/prometheus_alert_bot/pkg/core"
	"github.com/belousovalex/prometheus_alert_bot/pkg/parser"
	"github.com/belousovalex/prometheus_alert_bot/pkg/rocketchat"
	"github.com/belousovalex/prometheus_alert_bot/pkg/telegram"
)

func SenderFactory(senderConfig *parser.SenderConfig) core.Sender {
	switch senderConfig.Type {
	case "telegram":
		return telegram.MakeTelegramSender(
			senderConfig.TelegramBotToken,
			senderConfig.TelegramChatId,
		)
	case "rocketchat":
		return rocketchat.MakeRocketChatSender(
			senderConfig.RocketChatWebHook,
		)
	}
	panic(fmt.Sprintf("Cant make sender by %v", senderConfig))
}
