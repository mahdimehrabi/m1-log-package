package logger

import (
	"context"
	logv1 "github.com/mahdimehrabi/m1-log-proto/gen/go/log/v1"
	"github.com/rs/zerolog/log"
	"time"
)

const workerCount = 1000
const queueMaxLength = 20000
const grpcConnectionCount = 100

type sender struct {
	lsc           logv1.LogServiceClient //logger service client
	queue         chan message
	lStoreClients []logv1.LogService_StoreLogClient
}

func (s *sender) setup() error {
	s.lStoreClients = make([]logv1.LogService_StoreLogClient, grpcConnectionCount)
	for i := 0; i < grpcConnectionCount; i++ {
		lsc, err := s.lsc.StoreLog(context.Background())
		if err != nil {
			return err
		}
		s.lStoreClients[i] = lsc
	}
	s.runWorkers()
	return nil
}

func (s *sender) send(msg message) {
	s.queue <- msg
}

func (s *sender) runWorkers() {
	//dividing connections to workers
	lscI := len(s.lStoreClients)
	for i := 0; i < workerCount; i++ {
		lscI--
		if lscI < 0 {
			lscI = len(s.lStoreClients) - 1
		}
		go s.worker(lscI)
	}
}

func (s *sender) worker(lsci int) {
	for msg := range s.queue {
		if !msg.sentGrpc {
			err := s.lStoreClients[lsci].Send(&logv1.Log{
				Error: msg.message,
			})
			if err != nil {
				log.Err(err).Send()
			} else {
				msg.sentGrpc = true
			}
		}

		if !msg.sentLocal {
			switch msg.level {
			case errLevel:
				log.Error().Msg(msg.message)
			case infoLevel:
				log.Info().Msg(msg.message)
			case fatalLevel:
				log.Fatal().Msg(msg.message)
			case warningLevel:
				log.Warn().Msg(msg.message)
			}
			msg.sentLocal = true
		}
		if !msg.sentLocal || !msg.sentGrpc {
			//cool down grpc connection
			time.Sleep(100 * time.Millisecond)
			s.queue <- msg
		}
	}
}
func (s *sender) renewConnection(lsci int) {
	lsc, err := s.lsc.StoreLog(context.Background())
	if err != nil {
		log.Err(err).Send()
		time.Sleep(time.Second * 1)
		s.renewConnection(lsci)
		return
	}
	s.lStoreClients[lsci] = lsc
}
