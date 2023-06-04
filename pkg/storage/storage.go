package storage

import (
	"context"
)

// FilePart представляет собой часть файла
type FilePart struct {
	ID   string
	Part int
	Data []byte
}

// Storage представляет интерфейс для взаимодействия с серверами хранения
type Storage interface {
	Upload(ctx context.Context, fileParts []*FilePart) (err error)
	Download(ctx context.Context, id string) (fileParts []*FilePart, err error)
}
