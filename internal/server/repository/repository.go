package repository

import (
	"context"

	"github.com/mamadeusia/file-transfer-proof/internal/server/entity"
)

type CollectionRepository interface {
	StoreCollectionData(ctx context.Context, collectionHash string, fileData entity.FileData) error
	GetCollectionData(ctx context.Context, colletionHash string) ([]entity.FileData, error)
}

type MemoryRepository struct {
	l map[string][]entity.FileData
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		l: make(map[string][]entity.FileData),
	}
}

func (mr *MemoryRepository) StoreCollectionData(ctx context.Context, collectionHash string, fileData entity.FileData) error {
	mr.l[collectionHash] = append(mr.l[collectionHash], fileData)
	return nil
}
func (mr *MemoryRepository) GetCollectionData(ctx context.Context, colletionHash string) ([]entity.FileData, error) {
	return mr.l[colletionHash], nil
}
