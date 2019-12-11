package parser

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type CheckConfigError struct {
	messages []string
}

func (e CheckConfigError) Error() string {
	return strings.Join(e.messages, "\n")
}

type Patterns struct {
	Labels map[string]string `yaml:"labels"`
}

type SenderConfig struct {
	Type              string
	Patterns          *Patterns `yaml:"patterns"`
	TelegramBotToken  string    `yaml:"telegram_bot_token"`
	TelegramChatId    string    `yaml:"telegram_chat_id"`
	RocketChatWebHook string    `yaml:"rocketchat_webhook"`
}

type Config struct {
	Senders []*SenderConfig
}

func ParseConfigFile(pathToFile string) (*Config, error) {
	var err error
	var configData []byte
	var config Config
	configData, err = ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	errors := checkConfig(&config)
	if len(errors) > 0 {
		return nil, CheckConfigError{errors}
	}

	return &config, nil
}

func checkConfig(config *Config) []string {
	var errors []string

	for _, sender := range config.Senders {
		switch tp := sender.Type; tp {
		case "telegram":
			errors = append(errors, checkTelegramSender(sender)...)
		case "rocketchat":
			errors = append(errors, checkRocketChatSender(sender)...)
		default:
			errors = append(errors, fmt.Sprintf("Cant recognize %s type of messagers.", tp))
		}
	}
	return errors
}

func checkTelegramSender(sender *SenderConfig) []string {
	sender.TelegramChatId = resolveFromEnv(sender.TelegramChatId)
	sender.TelegramBotToken = resolveFromEnv(sender.TelegramBotToken)
	if len(sender.TelegramChatId) == 0 {
		return []string{"Parameter telegram_chat_id can not be empty."}
	}
	if len(sender.TelegramBotToken) == 0 {
		return []string{"Parameter telegram_bot_token can not be empty."}
	}
	return []string{}
}

func checkRocketChatSender(sender *SenderConfig) []string {
	sender.RocketChatWebHook = resolveFromEnv(sender.RocketChatWebHook)
	if len(sender.RocketChatWebHook) == 0 {
		return []string{"Parameter rocketchat_webhook can not be empty."}
	}
	return []string{}
}

func resolveFromEnv(originValue string) string {
	if strings.HasPrefix(originValue, "$") {
		nameOfEnvVar := string([]rune(originValue)[1:])
		envValue, ok := os.LookupEnv(nameOfEnvVar)
		if !ok {
			fmt.Printf("Environment %s does not set", nameOfEnvVar)
		}
		return envValue
	}
	return originValue
}
