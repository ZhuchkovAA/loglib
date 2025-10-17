package service

import (
	"context"
	consts "github.com/ZhuchkovAA/loglib/constants"
	"github.com/ZhuchkovAA/loglib/internal/domain/models"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type senderClient interface {
	WriteLog(context.Context, *LogRequest) (*LogResponse, error)
}

type Sender struct {
	client      senderClient
	ServiceName string
}

type LogRequest struct {
	Service   string
	Level     string
	Message   string
	Fields    map[string]*structpb.Value
	Timestamp int64
}

type LogResponse struct {
	Status      int32
	ProcessedAt *timestamppb.Timestamp
	LogId       string
}

func NewSender(serviceName string, client senderClient) *Sender {
	return &Sender{
		client:      client,
		ServiceName: serviceName,
	}
}

func (s *Sender) Send(log *models.LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := s.client.WriteLog(ctx, &LogRequest{
		Service:   s.ServiceName,
		Level:     consts.GetErrorLevelStr(log.Level),
		Message:   log.Message,
		Fields:    log.Metadata,
		Timestamp: log.Timestamp,
	})

	return err
}
