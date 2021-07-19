package main

import (
	"drone-delivery/drone-swarm/pkg/config"
	"drone-delivery/drone-swarm/pkg/domain/flying"
	"drone-delivery/drone-swarm/pkg/domain/routing"
	"drone-delivery/drone-swarm/pkg/domain/telemetry"
	"drone-delivery/drone-swarm/pkg/domain/warehouse"
	"drone-delivery/drone-swarm/pkg/network/inbound/http/rest"
	"drone-delivery/drone-swarm/pkg/network/outbound/http/grpc"
	"drone-delivery/drone-swarm/pkg/network/outbound/http/json"
	"fmt"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("ez a drón-raj a szimulacioban, a drónok példányait szimulálja.")
	config.SetConfig()
	var logger goKitLog.Logger
	logger = goKitLog.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo()) // <--
	logger = goKitLog.With(logger, "ts", goKitLog.DefaultTimestampUTC)
	jsonOutboundAdapter := json.NewOutBoundAdapter()
	grpcOutboundAdapter := grpc.NewOutBoundAdapter()

	telemetryService := telemetry.NewService(jsonOutboundAdapter, logger)

	geo := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	routingService := routing.NewService(geo)

	flyingService := flying.NewService(telemetryService, routingService, logger)
	warehouseService := warehouse.NewService(flyingService, logger)
	log.Fatal(http.ListenAndServe(":2000", rest.Handler(warehouseService, telemetryService, grpcOutboundAdapter, jsonOutboundAdapter)))
}

//ez csak egy sima kliens (drón a szimulacioban), ami megkapja a celt es ez alapjan fog az utvonal alatt mindenfele adatokat generalni es kuldeni magarol.
//Viszonylag buta, par dolgot kell jol csinalnia, az adatokat gyorsan es minimalis eroforras felhasznalasaval küldje a szervernek.
