package actions

import (
	"context"
	"github.com/ZhuchkovAA/loglib/internal/config"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"time"

	"github.com/ZhuchkovAA/protoRMRF/gen/go/log.v1"
	"google.golang.org/grpc"
)

type Sender struct {
	client  logv1.LogServiceClient
	Service string
}

func NewSender(cfg config.Config) *Sender {
	conn, err := grpc.Dial(cfg.GRPCAddress, grpc.WithInsecure())
	if err != nil {
		return nil
	}

	return &Sender{
		client:  logv1.NewLogServiceClient(conn),
		Service: cfg.ServiceName,
	}
}

func (s *Sender) Send(log models.LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := s.client.WriteLog(ctx, &logv1.LogRequest{
		Service:   log.Service,
		Level:     log.Level,
		Message:   log.Message,
		Metadata:  log.Metadata,
		Timestamp: log.Timestamp,
	})

	return err
}
