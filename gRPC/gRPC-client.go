package grpc

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client StorageClient
	conn   *grpc.ClientConn
	once   sync.Once
)

func Connect() StorageClient {
	once.Do(func() {
		var err error
		maxMsgSize := 1024 * 1024 * 2000 // 50 MB

		conn, err = grpc.NewClient("localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(maxMsgSize),
				grpc.MaxCallSendMsgSize(maxMsgSize),
			),
		)
		if err != nil {
			log.Fatalf("failed to connect to Python gRPC server at localhost:50051: %v", err)
		}
		client = NewStorageClient(conn)
	})
	return client
}

func Close() {
	if conn != nil {
		conn.Close()
	}
}
