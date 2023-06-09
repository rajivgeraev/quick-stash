package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
)

type BadgerStorageService struct {
	db *badger.DB
}

// NewBadgerStorageService создает новый экземпляр BadgerStorageService
func NewBadgerStorageService(db *badger.DB) *BadgerStorageService {
	return &BadgerStorageService{
		db: db,
	}
}

// StoreChunk сохраняет chunk файла в BadgerDB
func (s *BadgerStorageService) StoreChunk(id string, chunkNumber int, content io.Reader) error {
	data, err := ioutil.ReadAll(content)
	if err != nil {
		log.Printf("ReadAll error: %v", err)
		return err
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s_%d", id, chunkNumber)
		err := txn.Set([]byte(key), data)
		return err
	})

	if err != nil {
		log.Printf("BadgerDB Set error: %v", err)
		return err
	}

	return nil
}

// GetChunkNumbers возвращает список номеров чанков для данного id файла
func (s *BadgerStorageService) GetChunkNumbers(id string) ([]int, error) {
	var chunkNumbers []int
	err := s.db.View(func(txn *badger.Txn) error {
		prefix := []byte(id + "_")
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			keyParts := strings.Split(string(key), "_")
			if len(keyParts) < 2 {
				continue
			}
			chunkNumber, err := strconv.Atoi(keyParts[1])
			if err != nil {
				log.Printf("BadgerDB atoi error: %v", err)
				continue
			}
			chunkNumbers = append(chunkNumbers, chunkNumber)
		}
		return nil
	})

	if err != nil {
		log.Printf("BadgerDB View error: %v", err)
		return nil, err
	}

	return chunkNumbers, nil
}

// DeleteChunks удаляет все чанки для данного id файла
func (s *BadgerStorageService) DeleteChunks(id string) error {
	prefix := []byte(id + "_")
	return s.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			err := txn.Delete(key)
			if err != nil {
				log.Printf("BadgerDB Delete error: %v", err)
				return err
			}
		}

		return nil
	})
}

// GetChunk возвращает содержимое конкретного чанка файла.
func (s *BadgerStorageService) GetChunk(id string, chunkNumber int) (io.Reader, error) {
	key := []byte(fmt.Sprintf("%s_%d", id, chunkNumber))
	var valCopy []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(valCopy), nil
}
