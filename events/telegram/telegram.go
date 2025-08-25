package telegram

import (
	"errors"
	"tgBot/clients/telegram"
	"tgBot/events"
	e "tgBot/lib/error"
	"tgBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}

	p.offset = update[len(update)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ev events.Event) error {
	switch ev.Type {
	case events.Message:
		return p.processMessage(ev)
	default:
		return e.Wrap("can't process Message", errors.New("unknown event type"))
	}
}

func (p *Processor) processMessage(ev events.Event) error {
	m, err := meta(ev)
	if err != nil {
		return e.Wrap("can't get metadata", err)
	}

	if err := p.doCmd(ev.Text, m.ChatID, m.Username); err != nil {
		return e.Wrap("can't do cmd", err)
	}

	return nil
}

func meta(ev events.Event) (Meta, error) {
	if m, ok := ev.Meta.(Meta); ok {
		return m, nil
	}
	return Meta{}, e.Wrap("can't get meta from event", errors.New("can't get meta from event"))
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)

	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}
