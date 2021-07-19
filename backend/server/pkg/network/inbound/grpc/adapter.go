package grpc

import (
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/domain/services/telemetry"
	"drone-delivery/server/pkg/network/inbound/grpc/protobuf"
	"google.golang.org/grpc"
	"io"
	"log"
)

type adapter struct {
	telemetryService telemetry.Service
	server           *grpc.Server
	protobuf.UnimplementedTelemetryServiceServer
}

func NewAdapter(telemetryService telemetry.Service, server *grpc.Server) *adapter {
	return &adapter{telemetryService: telemetryService, server: server}
}

func (a *adapter) TelemetryStream(stream protobuf.TelemetryService_TelemetryStreamServer) error {
	log.Printf("SaveTelemetryStream function was invoked with a streaming request\n")
	result := "OK"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&protobuf.TelemetryDataResponse{Result: result})
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v\n", err)
			return err
		}
		t := req.GetTelemetry()
		var motorTemps []int
		var telemetryErrors []models.TelemetryError
		for i, v := range t.GetMotorTemperatures() {
			motorTemps[i] = int(v)
		}
		for i, v := range t.GetErrors() {
			telemetryErrors[i] = models.TelemetryError(int(v))
		}

		telemetryData := models.Telemetry{
			Speed: t.GetSpeed(),
			Location: models.GPS{
				Latitude:  t.GetLocation().GetLatitude(),
				Longitude: t.GetLocation().GetLongitude(),
			},
			Altitude:           t.GetAltitude(),
			Bearing:            t.GetBearing(),
			Acceleration:       t.GetAcceleration(),
			BatteryLevel:       int(t.GetBatteryLevel()),
			BatteryTemperature: int(t.GetBatteryTemperature()),
			MotorTemperatures:  motorTemps,
			Errors:             telemetryErrors,
			TimeStamp:          t.TimeStamp.AsTime(),
			DroneID:            int(t.GetDroneId()),
		}
		err = a.telemetryService.SaveTelemetry(telemetryData)
	}
}
