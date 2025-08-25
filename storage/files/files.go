package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	e "tgBot/lib/error"
	"tgBot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm = 0774
)

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(p *storage.Page) (err error) {

	filePath := filepath.Join(s.basePath, p.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return e.Wrap("can't create directory", err)
	}

	fileN, err := fileName(p)
	if err != nil {
		return e.Wrap("can't create file name", err)
	}

	filePath = filepath.Join(filePath, fileN)

	file, err := os.Create(filePath)
	if err != nil {
		return e.Wrap("can't create file", err)
	}

	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(p)
	if err != nil {
		return e.Wrap("can't encode page to file", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, e.Wrap("can't read directory", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]
	data, err := s.decodePage(file.Name())
	if err != nil {
		return nil, e.Wrap("can't decode file", err)
	}

	return data, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileN, err := fileName(p)
	if err != nil {
		return e.Wrap("can't get fileName", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileN)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprint("can't remove file: %s", path), err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileN, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't get fileName", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileN)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap("can't stat file", err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't open file", err)
	}

	defer func() { _ = f.Close() }()

	var page storage.Page
	if err = gob.NewDecoder(f).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}
