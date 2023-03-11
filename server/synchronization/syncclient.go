package synchronization

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
	"google.golang.org/grpc"
)

func (s *Sync) createClient(ctx context.Context, ratelimiterId int64, port int64) (*grpc.ClientConn, error) {
	// create a gRPC client for other rate limiter nodes
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Printf("[Sync] failed to connect to %d ERR: %v", port, err)
		return nil, err
	}
	// add conn to map for cleanup during shutdown
	s.NodeConns[port] = conn
	return conn, nil
}

func (s *Sync) createStream(ctx context.Context, conn *grpc.ClientConn, ratelimiterId int64, port int64) (Sync_SyncServiceClient, error) {
	// Create sync client stream
	cl := NewSyncClient(conn)
	stream, err := cl.SyncService(context.Background())
	if err != nil {
		log.Printf("[Sync] failed to create stream %d ERR: %v", port, err)
		return nil, err
	}
	s.NodeClients[ratelimiterId] = stream
	return stream, nil
}

// CLIENT
func (s *Sync) runClientReceiverAndSender(sync_ Sync_SyncServiceClient) {
	go s.clientReceiveAliveFromServer(sync_)
	go s.clientBroadcastSyncMsgToNodes(context.Background())
	// Apparently I also need this select {} here to keep process open...
	select {
	// Exit on stream context done
	case <-sync_.Context().Done():
		return
	}
}

// Client receiver
func (s *Sync) clientReceiveAliveFromServer(sync_ Sync_SyncServiceClient) {
	for {
		msg, err := sync_.Recv()
		if err != nil {
			if err == io.EOF {
				// Stream has been closed
				log.Println("[Sync] CLIENT receive stream closed")
				return
			}
			log.Printf("[Sync] CLIENT error when receiving msg. ERR: %v", err)
			return
		}
		if msg == nil {
			log.Println("[Sync] CLIENT received NIL msg")
		}
		if msg != nil {
			// Response || msg from SERVER will be ALIVE ping
			log.Printf("[Sync] CLIENT received incoming msg: %v", msg)
		}
	}
}

// client broadcaster/sender
func (s *Sync) clientBroadcastSyncMsgToNodes(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-s.rt.OutSyncChan:
			if !ok {
				log.Println("[Sync] CLIENT RateTracker out sync channel closed, unable to broadcast")
				return
			}
			log.Printf("[Sync] CLIENT sending boradcast of sync msg :: %v", out)
			go s.sendSyncMsg(ctx, out)
		default:
			// OutSyncChan is empty, wait a bit before checking channel
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (s *Sync) sendSyncMsg(ctx context.Context, out *ratetracker.SyncMsg) {
	for id, stream := range s.NodeClients {
		msg := SyncMsg{RatelimiterId: int64(*s.id), ApiKey: out.ApiKey, Timestamp: out.Timestamp}
		if err := stream.Send(&msg); err != nil {
			log.Printf("[Sync] CLIENT failed to send msg to :%d, maybe node is down ERR: %v", msg.RatelimiterId, err)
			s.deRegisterNode(id)
		}
	}
}

func (s *Sync) deRegisterNode(id int64) {
	// de-register nodes
	delete(s.NodeClients, id)
	delete(s.NodeConns, id)
}