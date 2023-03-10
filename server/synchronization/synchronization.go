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

type Sync struct {
	id			*int
	NodeClients map[int64]Sync_BroadcastServiceClient
	NodeConns   map[int64]*grpc.ClientConn
	rt          *ratetracker.RateTracker

	UnimplementedSyncServer
}

func NewSyncImpl(rt *ratetracker.RateTracker, port *int) *Sync {
	s := &Sync{
		NodeClients: make(map[int64]Sync_BroadcastServiceClient),
		NodeConns:	 make(map[int64]*grpc.ClientConn),
		rt:          rt,
		id:			 port,
	}

	return s
}

// SERVER
func (s *Sync) BroadcastService(sync_ Sync_BroadcastServiceServer) error {
	go s.receiveBroadcast(sync_)
	go s.SendBroadcast(sync_)
	select {}
}

// SERVER receiver
func (s *Sync) receiveBroadcast(sync_ Sync_BroadcastServiceServer) {
	for {
		msg, err := sync_.Recv()
		if err != nil {
			if err == io.EOF {
				// Stream has been closed
				log.Println("[Sync] SERVER stream closed")
				return
			}
			log.Printf("[Sync] SERVER error in broadcastService when receiving msg. ERR: %v", err)
			return
		}
		if msg == nil {
			log.Println("[Sync] SERVER received NIL msg")
		}
		if msg != nil {
			log.Printf("[Sync] SERVER received incoming msg: %v", msg)
			// Pushing incoming msg as events to InSyncChannel in RateTracker service
			// RateTracker service will handle these events coming in, and update ds
			syncMsg := ratetracker.SyncMsg{RatelimiterId: msg.RatelimiterId, ApiKey: msg.ApiKey, Timestamp: msg.Timestamp}
			s.rt.InSyncChan <- &syncMsg
		}
	}
}

// SERVER sender
func (s *Sync) SendBroadcast(sync_ Sync_BroadcastServiceServer) {
	for {
		// Start a ticker that executes each 2 seconds
		timer := time.NewTicker(2 * time.Second)
		select {
		// Exit on stream context done
		case <-sync_.Context().Done():
			return
		case <-timer.C:
			msg := SyncAck{Ack: true}
			if err := sync_.Send(&msg); err != nil {
				log.Printf("[Sync] SERVER failed to send ALIVE ERR: %v", err)
			}
		default:
			// OutSyncChan is empty, wait a bit before checking channel
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func (s *Sync) RegisterNode(ratelimiterId int64, port int64) error {
	// create connection to new node, and store conn and stream
	stream, err := s.createClient(port)
	if err != nil {
		log.Printf("[Sync] failed to create client while registering node ERR: %v", err)
		return err
	}
	log.Printf("[Sync] successfully registered stream client with :: %d", port)
	s.NodeClients[ratelimiterId] = stream
	// start goroutine for Client receiver/sender
	go s.pleaseWorkBitch(stream);

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

func (s *Sync) createClient(port int64) (Sync_BroadcastServiceClient, error) {
	// create a gRPC client for other rate limiter nodes
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Printf("[Sync] failed to connect to %d ERR: %v", port, err)
		return nil, err
	}
	// add conn to map for cleanup during shutdown
	s.NodeConns[port] = conn
	log.Printf("[Sync] created client connection with node :: %d", port)
	// Create sync client stream
	cl := NewSyncClient(conn)
	stream, err := cl.BroadcastService(context.Background())
	if err != nil {
		log.Printf("[Sync] failed to create stream %d ERR: %v", port, err)
		return nil, err
	}
	log.Printf("[Sync] created gRPC stream with :: %d", port)
	return stream, nil
}

// CLIENT
func (s *Sync) pleaseWorkBitch(sync_ Sync_BroadcastServiceClient) {
	go s.streamClientReceiver(sync_)
	go s.StreamClientSender()
	select {}
}

// Client receiver
func (s *Sync) streamClientReceiver(sync_ Sync_BroadcastServiceClient) {
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
func (s *Sync) StreamClientSender() {
	for {
		select {
		case val, ok := <-s.rt.OutSyncChan:
			if !ok {
				log.Println("[Sync] CLIENT RateTracker out sync channel closed, unable to broadcast")
				return
			}
			log.Printf("[Sync] CLIENT sending boradcast of sync msg :: %v", val)
			out := *val
			for _, stream := range s.NodeClients {
				msg := SyncMsg{RatelimiterId: int64(*s.id), ApiKey: out.ApiKey, Timestamp: out.Timestamp}
				if err := stream.Send(&msg); err != nil {
					log.Printf("[Sync] CLIENT failed to send msg to :%d, maybe node is down ERR: %v", msg.RatelimiterId, err)
				}
			}
		default:
			// OutSyncChan is empty, wait a bit before checking channel
			time.Sleep(time.Millisecond * 100)
		}
	}
}