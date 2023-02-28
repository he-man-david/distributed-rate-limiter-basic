package synchronization

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

type Sync struct {
	NodeClients map[int64]SyncClient

	UnimplementedSyncServer
}

func NewSyncImpl() *Sync {
	return &Sync{
		NodeClients: make(map[int64]SyncClient),
	}
}

func (s *Sync) BroadcastService(sync_ Sync_BroadcastServiceServer) error {
	s.receiveBroadcast(sync_)
	return nil
}

func (s *Sync) receiveBroadcast(sync_ Sync_BroadcastServiceServer) {
	for {
		msg, err := sync_.Recv()
		if err != nil {
			log.Printf("[Sync] error in broadcastService when receiving msg. ERR: %v", err)
		}
		log.Printf("[Sync] received incoming msg: %v", msg)
	}
}

func (s *Sync) SendBroadcast(ctx context.Context, sync_ Sync_BroadcastServiceServer, msg *SyncMsg) error {
	// TODO: David
	// Find all Node client
	err := sync_.SendMsg(&msg)
	if err != nil {
		log.Printf("[Sync] error in broadcastService when sending msg: %v \n ERR: %v", msg, err)
		return err
	}
	return nil
}

func (s *Sync) RegisterNode(ctx context.Context, ratelimiterId int64, port int64) error {
	newClient, err := createClient(port)
	if err != nil {
		return err
	}
	s.NodeClients[ratelimiterId] = newClient
	return nil
}

func createClient(port int64) (SyncClient, error) {
	// create a gRPC client for other rate limiter nodes
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Sync] failed to connect to %d: %v", port, err)
		return nil, err
	}
	return NewSyncClient(conn), nil
}