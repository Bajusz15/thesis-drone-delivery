package grpc

import (
	"context"
	"drone-delivery/drone-swarm/pkg/config"
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/network/inbound/grpc/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StreamClient struct {
	Tsc    protobuf.TelemetryServiceClient
	Stream protobuf.TelemetryService_TelemetryStreamClient
}

type Adapter struct {
	cc  *grpc.ClientConn
	tsc protobuf.TelemetryServiceClient
	//Stream protobuf.TelemetryService_TelemetryStreamClient
	streams map[int]StreamClient
}

func NewOutBoundAdapter() *Adapter {

	var err error
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	a := &Adapter{}
	a.streams = make(map[int]StreamClient)

	a.cc, err = grpc.Dial(config.ServerGRPCDomain+":"+config.ServerGRPCPort, opts...)
	if err != nil {
		panic(err)
	}
	//defer a.cc.Close()
	a.tsc = protobuf.NewTelemetryServiceClient(a.cc)
	//a.Stream, err = a.Tsc.TelemetryStream(context.TODO())
	//if err != nil {
	//	panic(err)
	//}
	return a
}

func (a *Adapter) SendTelemetryDataToServer(t models.Telemetry) error {
	var err error
	var streamer protobuf.TelemetryService_TelemetryStreamClient
	if sc, ok := a.streams[t.DroneID]; ok {
		//do something here
		streamer = sc.Stream
	} else {
		var tempStream = StreamClient{}
		tempStream.Tsc = protobuf.NewTelemetryServiceClient(a.cc)
		tempStream.Stream, err = tempStream.Tsc.TelemetryStream(context.TODO())
		if err != nil {
			return err
		}
		a.streams[t.DroneID] = tempStream
		streamer = a.streams[t.DroneID].Stream
	}

	temperatures := make([]int32, len(t.MotorTemperatures))
	telemetryErrors := make([]int32, len(t.Errors))
	for i, val := range t.MotorTemperatures {
		temperatures[i] = int32(val)
	}
	for i, val := range t.Errors {
		telemetryErrors[i] = int32(val)
	}
	telemetryDataRequest := protobuf.TelemetryDataRequest{
		Telemetry: &protobuf.Telemetry{
			Speed: t.Speed,
			Location: &protobuf.GPS{
				Latitude:  t.Location.Latitude,
				Longitude: t.Location.Longitude,
			},
			Altitude:           t.Altitude,
			Bearing:            t.Bearing,
			Acceleration:       t.Acceleration,
			BatteryLevel:       int32(t.BatteryLevel),
			BatteryTemperature: int32(t.BatteryTemperature),
			MotorTemperatures:  temperatures,
			Errors:             telemetryErrors,
			TimeStamp:          timestamppb.New(t.TimeStamp),
			DroneId:            int32(t.DroneID),
		},
	}
	err = streamer.Send(&telemetryDataRequest)
	return err
}
