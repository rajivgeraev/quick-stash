//pb "github.com/rajivgeraev/quick-stash/StorageServer/grpc/proto"

package grpc

import (
	"bytes"
	"context"
	"io"

	pb "github.com/rajivgeraev/quick-stash/StorageServer/grpc/proto"
)

type storageChunks interface {
	StoreChunk(id string, chunkNumber int, content io.Reader) error
	DeleteChunks(id string) error
	GetChunk(id string, chunkNumber int) (io.Reader, error)
	GetChunkNumbers(id string) ([]int, error)
}

// Server представляет собой реализацию сервера gRPC для обработки запросов на хранение.
type Server struct {
	pb.UnimplementedStorageServiceServer
	storage storageChunks
}

// NewServer создает новый объект сервера gRPC.
func NewServer(storage storageChunks) *Server {
	return &Server{
		storage: storage,
	}
}

// Store хранит чанк данных на сервере.
func (s *Server) Store(stream pb.StorageService_StoreServer) error {
	var (
		fileID      string
		chunkNumber int
	)

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StoreReply{
				Success: true,
				Message: "File stored successfully",
			})
		}
		if err != nil {
			return err
		}

		fileID = chunk.Id
		chunkNumber = int(chunk.ChunkNumber)
		content := bytes.NewReader(chunk.Content)
		err = s.storage.StoreChunk(fileID, chunkNumber, content)
		if err != nil {
			return err
		}
	}
}

// Delete удаляет все чанки данных с сервера.
func (s *Server) Delete(ctx context.Context, req *pb.FileRequest) (*pb.DeleteReply, error) {
	err := s.storage.DeleteChunks(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteReply{
		Success: true,
		Message: "File chunks deleted successfully",
	}, nil
}

// Retrieve возвращает чанк данных с сервера.
func (s *Server) Retrieve(req *pb.FileRequest, stream pb.StorageService_RetrieveServer) error {
	chunkNumbers, err := s.storage.GetChunkNumbers(req.Id)
	if err != nil {
		return err
	}

	for _, chunkNumber := range chunkNumbers {
		chunkContent, err := s.storage.GetChunk(req.Id, chunkNumber)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(chunkContent)

		chunk := &pb.FileChunk{
			Id:          req.Id,
			ChunkNumber: int32(chunkNumber),
			Content:     buf.Bytes(),
		}

		if err := stream.Send(chunk); err != nil {
			return err
		}
	}

	return nil
}
