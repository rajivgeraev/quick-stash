package grpcstorage

import (
	"context"
	"io"

	pb "github.com/rajivgeraev/quick-stash/api/proto"
	"github.com/rajivgeraev/quick-stash/pkg/storage"
	"google.golang.org/grpc"
)

type grpcStorage struct {
	client pb.FileStorageClient
}

func New(conn *grpc.ClientConn) storage.Storage {
	return &grpcStorage{client: pb.NewFileStorageClient(conn)}
}

func (s *grpcStorage) Upload(ctx context.Context, fileParts []*storage.FilePart) error {
	stream, err := s.client.Upload(ctx)
	if err != nil {
		return err
	}

	for _, part := range fileParts {
		if err := stream.Send(&pb.FileChunk{
			Id:    part.ID,
			Part:  int32(part.Part),
			Chunk: part.Data,
		}); err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}

func (s *grpcStorage) Download(ctx context.Context, id string) ([]*storage.FilePart, error) {
	stream, err := s.client.Download(ctx, &pb.DownloadRequest{Id: id})
	if err != nil {
		return nil, err
	}

	var fileParts []*storage.FilePart
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		fileParts = append(fileParts, &storage.FilePart{
			ID:   chunk.Id,
			Part: int(chunk.Part),
			Data: chunk.Chunk,
		})
	}

	return fileParts, nil
}
