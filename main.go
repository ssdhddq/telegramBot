package main

import (
	"flag"
	"log"
	"tgBot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	_ = telegram.New(tgBotHost, mustToken())

	//token = flags.Get(token)

	//tgClient = telegram.New(token)

	//fetcher = fetcher.New()

	//processor = processor.New()

	//consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
