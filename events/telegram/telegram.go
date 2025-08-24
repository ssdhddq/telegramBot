package telegram

import (
	"tgBot/clients/telegram"
	"tgBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client) {

}
