package sqlite

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	e "tgBot/lib/error"
	"tgBot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("sql don't open", err)
	}

	if err = db.Ping(); err != nil {
		return nil, e.Wrap("sql don't ping", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap("sql don't save page", err)
	}

	return nil
}
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("sql don't find page", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil

}
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap("sql don't remove page", err)
	}

	return nil
}
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count)
	if err != nil {
		return false, e.Wrap("sql don't find page", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.Wrap("sql don't create table", err)
	}

	return nil
}

func (s *Storage) GetAll(ctx context.Context, userName string) (*[]storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM()`

	var urls []string

	rows, err := s.db.QueryContext(ctx, q, userName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, e.Wrap("sql don't find page", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var pages = make([]storage.Page, 0, len(urls))
	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {
			return nil, e.Wrap("don't scan rows", err)
		}
		pages = append(pages, storage.Page{
			URL:      url,
			UserName: userName,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, e.Wrap("iterate rows", err)
	}

	if len(pages) == 0 {
		return nil, nil
	}

	return &pages, nil
}
