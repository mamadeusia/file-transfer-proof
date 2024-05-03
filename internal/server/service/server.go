package service

import (
	"bytes"
	"context"
	"io"

	config "github.com/mamadeusia/file-transfer-proof/config/server"
	"github.com/mamadeusia/file-transfer-proof/internal/server/entity"
	"github.com/mamadeusia/file-transfer-proof/internal/server/repository"
	"github.com/mamadeusia/file-transfer-proof/pkg/logger"
	"github.com/mamadeusia/file-transfer-proof/pkg/merkletree"
	uploadpb "github.com/mamadeusia/file-transfer-proof/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileServiceServer struct {
	uploadpb.UnimplementedFileServiceServer
	l    *logger.Logger
	cfg  *config.Config
	repo repository.CollectionRepository
}

func New(l *logger.Logger, cfg *config.Config, repo repository.CollectionRepository) *FileServiceServer {
	return &FileServiceServer{
		l:    l,
		cfg:  cfg,
		repo: repo,
	}
}

func (g *FileServiceServer) Upload(stream uploadpb.FileService_UploadServer) error {
	receivedFileIndex := make(chan int32, 10)
	defer close(receivedFileIndex)
	go func() {
		for {
			select {
			case rcv, ok := <-receivedFileIndex:
				if !ok {
					return
				}
				if err := stream.Send(&uploadpb.FileRecivedNotification{
					FileIndex: rcv,
				}); err != nil {
					// in the case of failure push it again in chan
					receivedFileIndex <- rcv
				}
			case <-stream.Context().Done():
				return
			}
		}
	}()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return g.logError(status.Error(codes.Internal, err.Error()))
		}
		file := NewFile()
		if err := file.SetFile(int(req.FileIndex), req.CollectionMerkleRoot, g.cfg.FilesStorage.Location); err != nil {
			return g.logError(status.Error(codes.Internal, err.Error()))
		}
		defer file.Close()
		reader := bytes.NewReader(req.GetContent())
		hash, err := merkletree.HashFromReader(reader)
		if err != nil {
			return g.logError(status.Error(codes.InvalidArgument, err.Error()))
		}
		if !merkletree.CheckProofWithLeafHash(req.GetCollectionMerkleRoot(), hash, req.GetProofs()) {
			return g.logError(status.Error(codes.InvalidArgument, "merkle root not matched"))
		}
		if _, err := reader.Seek(0, io.SeekStart); err != nil {
			return g.logError(status.Error(codes.Internal, err.Error()))
		}
		if err := file.WriteFromReader(reader); err != nil {
			return g.logError(status.Error(codes.Internal, err.Error()))
		}
		if err := g.repo.StoreCollectionData(stream.Context(), req.CollectionMerkleRoot, entity.FileData{
			Index:       req.FileIndex,
			ContentHash: hash,
		}); err != nil {
			return g.logError(status.Error(codes.Internal, err.Error()))
		}
		receivedFileIndex <- req.FileIndex
	}

	return nil
}

func (g *FileServiceServer) Download(ctx context.Context, req *uploadpb.DownloadRequest) (*uploadpb.DownloadResponse, error) {
	fd, err := g.repo.GetCollectionData(ctx, req.CollectionMerkleRoot)
	if err != nil {
		return nil, err
	}
	entity.SortFileData(fd)
	if len(fd) < int(req.FileIndex) {
		return nil, g.logError(status.Error(codes.OutOfRange, "totally not received"))
	}

	leafs := entity.FileDataContentHashes(fd)
	if merkletree.RootToString(merkletree.BuildTreeWithLeafHashes(leafs)) != req.CollectionMerkleRoot {
		return nil, g.logError(status.Error(codes.OutOfRange, "totally not received"))
	}

	file := NewFile()
	defer file.Close()
	if err := file.GetFile(int(req.FileIndex), req.CollectionMerkleRoot, g.cfg.FilesStorage.Location); err != nil {
		return nil, g.logError(status.Error(codes.Internal, err.Error()))
	}

	rsp := &uploadpb.DownloadResponse{
		Proofs: merkletree.ProofToStringSlice(merkletree.GetProofIndexWithLeafHashes(int(req.FileIndex), leafs)),
	}
	b, err := file.GetBytes()
	if err != nil {
		return nil, err
	}
	rsp.Content = b

	return rsp, nil
}
