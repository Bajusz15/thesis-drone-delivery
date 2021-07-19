package telemetry

import (
	"drone-delivery/server/pkg/domain/models"
	goKitLog "github.com/go-kit/kit/log"
)

type Service interface {
	SendTelemetry(t models.Telemetry) error
	ChangeService(a OutboundAdapter)
}

type OutboundAdapter interface {
	SendTelemetryDataToServer(t models.Telemetry) error
}

type service struct {
	adapter OutboundAdapter
	logger  goKitLog.Logger
}

func (s service) SendTelemetry(t models.Telemetry) error {
	err := s.adapter.SendTelemetryDataToServer(t)
	if err != nil {
		s.logger.Log("err", err, "desc", "outbound adapter returned with error")
		return err
	}
	return nil
}

func NewService(a OutboundAdapter, l goKitLog.Logger) *service {
	return &service{a, l}
}

func (s *service) ChangeService(a OutboundAdapter) {
	s.adapter = a
}
