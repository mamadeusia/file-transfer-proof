package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mamadeusia/file-transfer-proof/pkg/merkletree"
	uploadpb "github.com/mamadeusia/file-transfer-proof/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientService struct {
	addr         string
	diretoryPath string
	client       uploadpb.FileServiceClient
}

func New(addr string, diretoryPath string) *ClientService {
	return &ClientService{
		addr:         addr,
		diretoryPath: diretoryPath,
	}
}

func (s *ClientService) SendFile() error {

	conn, err := grpc.Dial(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	s.client = uploadpb.NewFileServiceClient(conn)
	interrupt := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(s *ClientService) {
		if err = s.upload(ctx, cancel); err != nil {
			log.Fatal(err)
			cancel()
		}
	}(s)

	select {
	case killSignal := <-interrupt:
		log.Println("Got ", killSignal)
		cancel()
	case <-ctx.Done():
	}
	return nil
}

func (s *ClientService) upload(ctx context.Context, cancel context.CancelFunc) error {
	stream, err := s.client.Upload(ctx)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(s.diretoryPath)
	if err != nil {
		return err
	}
	numFiles := 0
	var fileHashed []string
	for _, item := range files {
		if item.IsDir() {
			continue
		}
		file, err := os.Open(filepath.Join(s.diretoryPath, item.Name()))
		if err != nil {
			return err
		}
		h, err := merkletree.HashFromReader(file)
		if err != nil {
			return err
		}
		fileHashed = append(fileHashed, h)
		numFiles++
	}

	collectionHash := merkletree.RootToString(merkletree.BuildTreeWithLeafHashes(fileHashed))
	// TODO :: make concurrent streams send files
	// itemNameCh := make(chan string, 5)
	// for _, item := range files {
	// 	itemNameCh <- item.Name()
	// }
	errCh := make(chan error)
	go func(errCh chan error) {
		for i, item := range files {
			fileContentBytes, err := os.ReadFile(filepath.Join(s.diretoryPath, item.Name()))
			if err != nil {
				errCh <- err
			}
			err = stream.Send(&uploadpb.FileUploadRequest{
				FileIndex:            int32(i),
				Content:              fileContentBytes,
				CollectionMerkleRoot: collectionHash,
				Proofs:               merkletree.ProofToStringSlice(merkletree.GetProofIndexWithLeafHashes(i, fileHashed)),
			})
			if err != nil {
				errCh <- err
			}
		}
	}(errCh)

	receivedFileCnt := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		receivedFileCnt++
		fmt.Println(receivedFileCnt, " , ", numFiles)
		if receivedFileCnt == numFiles {
			cancel()
		}
	}

	return nil
}
