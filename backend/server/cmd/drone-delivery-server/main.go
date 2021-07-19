package main

import (
	"drone-delivery/server/pkg/config"
	"drone-delivery/server/pkg/domain/services/drone"
	"drone-delivery/server/pkg/domain/services/routing"
	"drone-delivery/server/pkg/domain/services/telemetry"
	grpcInbound "drone-delivery/server/pkg/network/inbound/grpc"
	"drone-delivery/server/pkg/network/inbound/grpc/protobuf"
	"drone-delivery/server/pkg/network/inbound/http/rest"
	"drone-delivery/server/pkg/network/outbound"
	"drone-delivery/server/pkg/storage/mongodb"
	"drone-delivery/server/pkg/storage/postgres"
	"fmt"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {
	fmt.Println("ez lesz a szerver, az adatközpont a szimulációban")
	var err error
	config.SetConfig()
	//logger
	var logger goKitLog.Logger
	logger = goKitLog.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo()) // <--
	logger = goKitLog.With(logger, "ts", goKitLog.DefaultTimestampUTC)

	//storage
	postgresStorage, err := postgres.NewStorage(config.PostgresConfig)
	if err != nil {
		log.Println("Error connecting to database")
		panic(err)
	}
	mongoStorage, err := mongodb.NewStorage(config.MongoConfig)
	if err != nil {
		log.Println("Error connecting to database")
		panic(err)
	}

	//err = generateDeliveryData()
	//if err != nil {
	//	logger.Log("err", err, "desc", "failed to generate delivery data")
	//	panic(err)
	//}
	//outbound adapters
	jsonAdapter := outbound.NewJSONAdapter()

	geo := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	//services
	var ts telemetry.Service
	var ds drone.Service
	var rs routing.Service
	ts = telemetry.NewService(postgresStorage, logger, geo)
	rs = routing.NewService(logger)
	ds = drone.NewService(postgresStorage, jsonAdapter, logger, rs)
	//REST API
	router := rest.Handler(ds, ts, postgresStorage, mongoStorage)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	//http
	go func(router http.Handler) {
		log.Println("serving http on port: 5000...")
		log.Fatal(http.ListenAndServe(":5000", router))
		wg.Done()
	}(router)

	//gRPC
	go func(service *telemetry.Service) {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal("Error creating to tcp listener")
		}
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		grpcAdapter := grpcInbound.NewAdapter(ts, grpcServer)
		protobuf.RegisterTelemetryServiceServer(grpcServer, grpcAdapter)
		log.Println("serving grpc on port 50051")
		log.Fatal(grpcServer.Serve(listener))
		wg.Done()
	}(&ts)

	wg.Wait()

}
