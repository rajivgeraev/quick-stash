package service

import "io"

type StorageService interface {
	StoreChunk(id string, chunkNumber int, content io.Writer) error
	DeleteChunks(id string) error
	GetChunk(id string, chunkNumber int) (io.Reader, error)
	GetChunkNumbers(id string) ([]int, error)
}
