package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/mamadeusia/file-transfer-proof/pkg/merkletree"
	uploadpb "github.com/mamadeusia/file-transfer-proof/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *ClientService) DownloadFile(index int32, collectionHash string) error {

	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	s.client = uploadpb.NewFileServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rsp, err := s.client.Download(ctx, &uploadpb.DownloadRequest{
		CollectionMerkleRoot: collectionHash,
		FileIndex:            index,
	})
	if err != nil {
		return err
	}
	contentHash, err := merkletree.HashFromReader(bytes.NewReader(rsp.Content))
	if err != nil {
		return err
	}
	if !merkletree.CheckProofWithLeafHash(collectionHash, contentHash, rsp.Proofs) {
		return errors.New("not matched")
	}
	fmt.Println("content is true")

	return nil
}
