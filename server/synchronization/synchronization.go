package synchronization

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/he-man-david/distributed-rate-limiter-basic/server/ratetracker"
	"google.golang.org/grpc"
)

type Sync struct {
	id			*int
	NodeClients map[int64]Sync_SyncServiceClient
	NodeConns   map[int64]*grpc.ClientConn
	rt          *ratetracker.RateTracker

	UnimplementedSyncServer
}

func NewSyncImpl(rt *ratetracker.RateTracker, port *int) *Sync {
	s := &Sync{
		NodeClients: make(map[int64]Sync_SyncServiceClient),
		NodeConns:	 make(map[int64]*grpc.ClientConn),
		rt:          rt,
		id:			 port,
	}

	return s
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

// SERVER
func (s *Sync) SyncService(sync_ Sync_SyncServiceServer) error {
	go s.receiveSyncFromClients(sync_)
	go s.sendAliveToClients(sync_)

	// Apparently I also need this select {} here to keep process open...
	select {
	// Exit on stream context done
	case <-sync_.Context().Done():
		return nil
	}
}

// SERVER receiver
func (s *Sync) receiveSyncFromClients(sync_ Sync_SyncServiceServer) error {
	for {
		select {
		// Exit on stream context done
		case <-sync_.Context().Done():
			return nil
		default:
			msg, err := sync_.Recv()
			if err != nil {
				if err == io.EOF {
					// Stream has been closed
					log.Println("[Sync] SERVER stream closed")
					return err
				}
				log.Printf("[Sync] SERVER error in broadcastService when receiving msg. ERR: %v", err)
				return err
			}
			if msg == nil {
				e := "[Sync] SERVER received NIL msg"
				log.Println(e)
				return errors.New(e)
			}
			log.Printf("[Sync] SERVER received incoming msg: %v", msg)
			// Pushing incoming msg as events to InSyncChannel in RateTracker service
			// RateTracker service will handle these events coming in, and update ds
			syncMsg := ratetracker.SyncMsg{RatelimiterId: msg.RatelimiterId, ApiKey: msg.ApiKey, Timestamp: msg.Timestamp}
			s.rt.InSyncChan <- &syncMsg
		}
	}
}

// SERVER sender (ALIVE msg)
func (s *Sync) sendAliveToClients(sync_ Sync_SyncServiceServer) error {
	timer := time.NewTicker(5 * time.Second)

	for {
		select {
			// Exit on stream context done
			case <-sync_.Context().Done():
				return nil
			case <-timer.C:
				// Start a ticker that executes every 5 seconds
				msg := Alive{Res: true}
				if err := sync_.Send(&msg); err != nil {
					log.Printf("[Sync] SERVER failed to send ALIVE ERR: %v", err)
					return err
				}
			default:
				// no need to check all the time
				time.Sleep(time.Millisecond * 500)
		}
	}
}

// When new nodes are started, and need to register its client
func (s *Sync) RegisterNode(ctx context.Context, ratelimiterId int64, port int64) error {
	// create conns with node server, and store conns
	conn, err := s.createClient(ctx, ratelimiterId, port)
	if err != nil {
		log.Printf("[Sync] failed to create client connection while registering node ERR: %v", err)
		return err
	}
	// create gRPC stream
	stream, err := s.createStream(ctx, conn, ratelimiterId, port)
	if err != nil {
		log.Printf("[Sync] failed to create stream while registering node ERR: %v", err)
		return err
	}
	log.Printf("[Sync] successfully registered stream client with :: %d", ratelimiterId)
	// start goroutine for Client receiver/sender
	go s.runClientReceiverAndSender(stream)

	return nil
}
