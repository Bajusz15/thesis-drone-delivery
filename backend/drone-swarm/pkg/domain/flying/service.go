package flying

import (
	"drone-delivery/drone-swarm/pkg/domain/routing"
	"drone-delivery/drone-swarm/pkg/domain/telemetry"
	"drone-delivery/drone-swarm/pkg/domain/warehouse"
	"drone-delivery/server/pkg/domain/models"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"log"
	"time"
)

type Service interface {
	StartFlight(d warehouse.Drone) error
	DropParcel() error
}

type service struct {
	//repo   Repository
	ts     telemetry.Service
	rs     routing.Service
	logger goKitLog.Logger
}

func NewService(ts telemetry.Service, rs routing.Service, l goKitLog.Logger) Service {
	return &service{ts, rs, l}
}

func (s *service) StartFlight(d warehouse.Drone) error {
	s.logger.Log("drone_id", d.ID, "desc", "drone started flying, ready to send telemetry")
	//TODO: ezt szakdolgozatba esetleg bele lehet irni mint erdekesseg, ahogyan meg van oldva ez a szimulacio
	period := time.Duration(200) * time.Millisecond
	ticker := time.NewTicker(period)
	//var wg sync.WaitGroup
	//wg.Add(1)
	//in parameter: group *sync.WaitGroup,
	go func(drone *warehouse.Drone, logger goKitLog.Logger) {
		//defer wg.Done()
		i := 0
		for range ticker.C {
			var err error
			destination := drone.Destinations[i]
			distance, bearing := s.rs.CalculateDroneDistanceAndDirectionFromDestination(
				drone.LastTelemetry.Location.Latitude,
				drone.LastTelemetry.Location.Longitude,
				destination.Coordinates.Latitude,
				destination.Coordinates.Longitude,
			)

			//distanceTravelled := drone.LastTelemetry.Speed * float64(period/(time.Duration(1)*time.Second)) //get meters travelled in a second)
			distanceTravelled := drone.LastTelemetry.Speed * 0.2 //get meters travelled in a second)
			currentLat, currentLon := s.rs.CalculateDroneNextCoordinates(
				drone.LastTelemetry.Location.Latitude,
				drone.LastTelemetry.Location.Longitude,
				distanceTravelled,
				bearing,
			)

			t := models.Telemetry{
				Speed: drone.LastTelemetry.Speed,
				Location: models.GPS{
					Latitude:  currentLat,
					Longitude: currentLon,
				},
				Altitude:           50,
				Bearing:            bearing,
				Acceleration:       0,
				BatteryLevel:       drone.LastTelemetry.BatteryLevel,
				BatteryTemperature: 0,
				MotorTemperatures:  nil,
				Errors:             nil,
				TimeStamp:          time.Now(), // <- ticker.C
				DroneID:            drone.ID,
			}

			if distance < 150 && t.Speed > 2 {
				t.Speed -= 1
			} else if distance < 12 {
				t.Speed = 0.2
				if distance < 1 {
					t.Speed = 0
					t.Location = drone.Destinations[i].Coordinates
					if drone.Destinations[i].WarehouseDestination {
						logger.Log("desc", "the drone successfully delivered the parcel and got back to the warehouse")
						return
					}
					if i < len(drone.Destinations) {
						i++
					}
					//move to next destination
					err = s.DropParcel()
					if err != nil {
						level.Warn(logger).Log("drone_id", drone.ID, "err", err, "desc", "failed to drop parcel")
						t.Errors = append(t.Errors, models.FailedToEjectPackage)
					}
				}
			} else if distance > 150 && t.Speed < 12 {
				t.Speed += 2
			}

			drone.LastTelemetry = t
			log.Println(drone.LastTelemetry)
			err = s.ts.SendTelemetry(t)
			if err != nil {
				level.Warn(logger).Log("err", err, "desc", "failed to send telemetry data")
			}

			//if signal is weak, change data send interval
			//ticker.Reset(time.Duration(500) * time.Millisecond)

			//if done, call ticker.Stop()
		}
		level.Info(logger).Log("drone_id", drone.ID, "desc", "drone stopped")
	}(&d, s.logger) // &wg,
	//wg.Wait()
	//main logic of flying return
	return nil
}

func (s *service) DropParcel() error {
	// 0.1 % chance of dropping the package
	//min := 1
	//max := 1000
	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	//number := r1.Intn(max-min+1) + min
	//if number == 1 {
	//	return ErrFailedToEjectPackage
	//}
	return nil
}

//StartFlight alternativ logika go.15 elott
//ticker := time.NewTicker(200 * time.Microsecond)
//quit := make(chan struct{})
//go func() {
//	for {
//		select {
//		case <- ticker.C:
//			// do stuff
//			ticker = time.NewTicker(500 * time.Microsecond)
//			//if done, quit the channel
//			//close(quit)
//		case <- quit:
//			ticker.Stop()
//			return
//		}
//	}
//}()
