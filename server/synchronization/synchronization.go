package synchronization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
	"google.golang.org/grpc"
)

type Sync struct {
	NodeClients map[int64]Sync_BroadcastServiceClient
	NodeConns   map[int64]*grpc.ClientConn
	rt          *ratetracker.RateTracker

	UnimplementedSyncServer
}

func NewSyncImpl(rt *ratetracker.RateTracker) *Sync {
	return &Sync{
		NodeClients: make(map[int64]Sync_BroadcastServiceClient),
		rt:          rt,
	}
}

func (s *Sync) BroadcastService(sync_ Sync_BroadcastServiceServer) error {
	go s.receiveBroadcast(sync_)
	return nil
}

func (s *Sync) receiveBroadcast(sync_ Sync_BroadcastServiceServer) {
	for {
		msg, err := sync_.Recv()
		if err != nil {
			log.Printf("[Sync] error in broadcastService when receiving msg. ERR: %v", err)
		}
		log.Printf("[Sync] received incoming msg: %v", msg)
		// Pushing incoming msg as events to InSyncChannel in RateTracker service
		// RateTracker service will handle these events coming in, and update ds
		syncMsg := ratetracker.SyncMsg{RatelimiterId: msg.RatelimiterId, ApiKey: msg.ApiKey, Timestamp: msg.Timestamp}
		s.rt.InSyncChan <- &syncMsg
	}
}

func (s *Sync) SendBroadcast(ctx context.Context, sync_ Sync_BroadcastServiceServer, msg *SyncMsg) {
	for {
		select {
		case val, ok := <-s.rt.OutSyncChan:
			if !ok {
				log.Println("[Sync] RateTracker out sync channel closed, unable to broadcast")
				return
			}
			for id, client := range s.NodeClients {
				msg := SyncMsg{RatelimiterId: id, ApiKey: val.ApiKey, Timestamp: val.Timestamp}
				go s.SendMsg(ctx, client, &msg)
			}
		default:
			// OutSyncChan is empty, wait a bit before checking channel
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (s *Sync) SendMsg(ctx context.Context, client Sync_BroadcastServiceClient, msg *SyncMsg) {
	if err := client.Send(msg); err != nil {
		log.Printf("[Sync] failed to broadcast msg to :%d, maybe node is down ERR: %v", msg.RatelimiterId, err)
	}
}

func (s *Sync) RegisterNode(ctx context.Context, ratelimiterId int64, port int64) error {
	// create connection to new node, and store conn and stream
	newClient, err := s.createClient(ctx, port)
	if err != nil {
		return err
	}
	s.NodeClients[ratelimiterId] = newClient
	return nil
}

// shutdown cleanup of all open connections and streams
func (s *Sync) Shutdown() {
	for id, client := range s.NodeClients {
		log.Printf("[Sync] shutting down node stream: %d ...", id)
		if err := client.CloseSend(); err != nil {
			log.Printf("[Sync] failed to shudown node stream %d ERR: %v", id, err)
		}
	}
	for id, conn := range s.NodeConns {
		log.Printf("[Sync] shutting down node conns: %d ...", id)
		if err := conn.Close(); err != nil {
			log.Printf("[Sync] failed to shudown node conns %d ERR: %v", id, err)
		}
	}
}

func (s *Sync) createClient(ctx context.Context, port int64) (Sync_BroadcastServiceClient, error) {
	// create a gRPC client for other rate limiter nodes
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Printf("[Sync] failed to connect to %d ERR: %v", port, err)
		return nil, err
	}
	// add conn to map for cleanup during shutdown
	s.NodeConns[port] = conn

	// Create sync client stream
	cl := NewSyncClient(conn)
	stream, err := cl.BroadcastService(ctx)
	if err != nil {
		log.Printf("[Sync] failed to create stream %d ERR: %v", port, err)
		return nil, err
	}

	return stream, nil
}
