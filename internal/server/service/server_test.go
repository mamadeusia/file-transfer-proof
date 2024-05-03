package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"log"
	"net"
	"os"
	"testing"
	"time"

	config "github.com/mamadeusia/file-transfer-proof/config/server"
	"github.com/mamadeusia/file-transfer-proof/internal/server/repository"
	"github.com/mamadeusia/file-transfer-proof/pkg/merkletree"
	uploadpb "github.com/mamadeusia/file-transfer-proof/pkg/proto"
	"github.com/stretchr/testify/assert"

	"github.com/mamadeusia/file-transfer-proof/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func newGrpcConn(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(20 * 1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func newClientStream(t *testing.T) uploadpb.FileService_UploadClient {
	l := logger.New("")
	serverPath := "/tmp/test_folder"

	defer os.RemoveAll(serverPath)
	c := config.Config{
		FilesStorage: config.FilesStorage{
			Location: serverPath,
		},
	}

	fileServer := New(l, &c, repository.NewMemoryRepository())

	conn := newGrpcConn(t, func(srv *grpc.Server) {
		uploadpb.RegisterFileServiceServer(srv, fileServer)
	})
	client := uploadpb.NewFileServiceClient(conn)
	context := context.Background()
	stream, err := client.Upload(context)
	if err != nil {
		t.Fatal(err)
	}
	return stream
}

func TestService_Upload_Download(t *testing.T) {

	client := newClientStream(t)

	numFiles := 10

	fileSize := 1024 * 1024

	var fileHashed []string
	var fileContents [][]byte

	// Create multiple files with random data
	for i := 0; i < numFiles; i++ {
		randData := make([]byte, fileSize)

		rand.Read(randData)
		h, err := merkletree.HashFromReader(bytes.NewReader(randData))
		if err != nil {
			t.Fatal(err)
		}
		fileHashed = append(fileHashed, h)
		fileContents = append(fileContents, randData)

	}

	collectionHash := merkletree.RootToString(merkletree.BuildTreeWithLeafHashes(fileHashed))
	// bCh := make(chan []byte)
	go func() {
		for i := 0; i < numFiles; i++ {
			err := client.Send(&uploadpb.FileUploadRequest{
				FileIndex:            int32(i),
				Content:              fileContents[i],
				CollectionMerkleRoot: collectionHash,
				Proofs:               merkletree.ProofToStringSlice(merkletree.GetProofIndexWithLeafHashes(i, fileHashed)),
			})
			if err != nil {
				t.Fatal(err)
			}
		}
	}()
	receivedFileCnt := 0
	for {
		_, err := client.Recv()
		if err != nil {
			break
		}
		receivedFileCnt++
		if receivedFileCnt == numFiles {
			err = client.CloseSend()
			assert.Nil(t, err)
		}
	}
	assert.Equal(t, receivedFileCnt, numFiles)

}
