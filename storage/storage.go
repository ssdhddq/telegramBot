package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	e "tgBot/lib/error"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't create hash on URL", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't create hash on UserName", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
