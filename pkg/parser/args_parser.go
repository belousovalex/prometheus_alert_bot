package parser

import "flag"

type BotFlags struct {
	ConfigPath string
}

func GetFlags() *BotFlags {
	var flags BotFlags

	flag.StringVar(
		&flags.ConfigPath,
		"config",
		"/etc/prometheus_alert_bot/config.yml",
		"Put path to bot config.",
	)
	flag.Parse()

	return &flags
}
