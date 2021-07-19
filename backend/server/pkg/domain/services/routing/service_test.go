package routing

import (
	"drone-delivery/server/pkg/domain/models"
	"fmt"
	goKitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestOptimizeRoutes(t *testing.T) {
	d1, d2, d3, d4 := models.Drone{
		ID:          1,
		Telemetry:   models.Telemetry{},
		Parcel:      models.Parcel{},
		Consumption: 800,
	}, models.Drone{
		ID:          2,
		Telemetry:   models.Telemetry{},
		Parcel:      models.Parcel{},
		Consumption: 700,
	}, models.Drone{
		ID:          3,
		Telemetry:   models.Telemetry{},
		Parcel:      models.Parcel{},
		Consumption: 650,
	}, models.Drone{
		ID:          4,
		Telemetry:   models.Telemetry{},
		Parcel:      models.Parcel{},
		Consumption: 390,
	}

	p1, p2, p3, p4 := models.Parcel{
		ID:     1,
		Weight: 0,
		DropOffSite: models.GPS{
			Latitude:  48.0745806,
			Longitude: 20.7696448,
		},
	}, models.Parcel{
		ID:     2,
		Weight: 0,
		DropOffSite: models.GPS{
			Latitude:  48.0736485,
			Longitude: 20.7791302,
		},
	}, models.Parcel{
		ID:     3,
		Weight: 0,
		DropOffSite: models.GPS{
			Latitude:  48.0637455,
			Longitude: 20.7497039,
		},
	}, models.Parcel{
		ID:     4,
		Weight: 0,
		DropOffSite: models.GPS{
			Latitude:  48.066926,
			Longitude: 20.7528756,
		},
	}

	var wh = models.Warehouse{Location: models.GPS{
		Latitude:  48.080922,
		Longitude: 20.766208,
	}}
	drones := []models.Drone{
		d1, d2, d3, d4,
	}
	parcels := []models.Parcel{
		p1, p2, p3, p4,
	}
	var logger goKitLog.Logger
	logger = goKitLog.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo()) // <--
	logger = goKitLog.With(logger, "ts", goKitLog.DefaultTimestampUTC)
	s := NewService(logger)
	fmt.Printf("distance is %v km", s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p1.DropOffSite.Latitude, p1.DropOffSite.Longitude, "K"))
	fmt.Printf("distance is %v km", s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p2.DropOffSite.Latitude, p2.DropOffSite.Longitude, "K"))
	fmt.Printf("distance is %v km", s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p3.DropOffSite.Latitude, p3.DropOffSite.Longitude, "K"))
	fmt.Printf("distance is %v km", s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p4.DropOffSite.Latitude, p4.DropOffSite.Longitude, "K"))
	log.Println("costs of parcels with drone 1")
	cost11 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p1.DropOffSite.Latitude, p1.DropOffSite.Longitude, "K") * d1.Consumption
	cost12 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p2.DropOffSite.Latitude, p2.DropOffSite.Longitude, "K") * d1.Consumption
	cost13 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p3.DropOffSite.Latitude, p3.DropOffSite.Longitude, "K") * d1.Consumption
	cost14 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p4.DropOffSite.Latitude, p4.DropOffSite.Longitude, "K") * d1.Consumption
	log.Println(cost11, cost12, cost13, cost14)
	log.Println("costs of parcels with drone 2")
	cost21 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p1.DropOffSite.Latitude, p1.DropOffSite.Longitude, "K") * d2.Consumption
	cost22 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p2.DropOffSite.Latitude, p2.DropOffSite.Longitude, "K") * d2.Consumption
	cost23 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p3.DropOffSite.Latitude, p3.DropOffSite.Longitude, "K") * d2.Consumption
	cost24 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p4.DropOffSite.Latitude, p4.DropOffSite.Longitude, "K") * d2.Consumption
	log.Println(cost21, cost22, cost23, cost24)
	log.Println("costs of parcels with drone 3")
	cost31 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p1.DropOffSite.Latitude, p1.DropOffSite.Longitude, "K") * d3.Consumption
	cost32 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p2.DropOffSite.Latitude, p2.DropOffSite.Longitude, "K") * d3.Consumption
	cost33 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p3.DropOffSite.Latitude, p3.DropOffSite.Longitude, "K") * d3.Consumption
	cost34 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p4.DropOffSite.Latitude, p4.DropOffSite.Longitude, "K") * d3.Consumption
	log.Println(cost31, cost32, cost33, cost34)
	log.Println("costs of parcels with drone 4")
	cost41 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p1.DropOffSite.Latitude, p1.DropOffSite.Longitude, "K") * d4.Consumption
	cost42 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p2.DropOffSite.Latitude, p2.DropOffSite.Longitude, "K") * d4.Consumption
	cost43 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p3.DropOffSite.Latitude, p3.DropOffSite.Longitude, "K") * d4.Consumption
	cost44 := s.CalculateDistance(wh.Location.Latitude, wh.Location.Longitude, p4.DropOffSite.Latitude, p4.DropOffSite.Longitude, "K") * d4.Consumption
	log.Println(cost41, cost42, cost43, cost44)
	dronesWithRoutes, err := s.OptimizeRoutes(wh, drones, parcels)
	if err != nil {
		s.logger.Log("desc", "failed to set routes for drones")
	}
	//after solving the assigment problem, these are the original costs, with the solution marked with \
	//\599.92	 524.93		487.43  	292.46
	//1004.19 	\878.67		815.9		489.54
	//1815.68	1588.72		1475.24		\885.14
	//1475.76   1291.29		\1199.1		719.43
	//The optimal value equals 3562.83
	//https://www.hungarianalgorithm.com/solve.php?c=599.92-524.93-487.43-292.46--1004.19-878.67-815.9-489.54--1815.68-1588.72-1475.24-885.14--1475.76-1291.29-1199.1-719.43&random=1

	assert.Equal(t, 1, dronesWithRoutes[0].Parcel.ID)
	assert.Equal(t, 2, dronesWithRoutes[1].Parcel.ID)
	assert.Equal(t, 4, dronesWithRoutes[2].Parcel.ID)
	assert.Equal(t, 3, dronesWithRoutes[3].Parcel.ID)

}

func TestCalculateDistance(t *testing.T) {
	var logger goKitLog.Logger
	logger = goKitLog.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo()) // <--
	logger = goKitLog.With(logger, "ts", goKitLog.DefaultTimestampUTC)
	s := NewService(logger)
	dist := s.CalculateDistance(48.080922, 20.766208, 48.0736485, 20.7791302)
	fmt.Printf("distance is %v km", dist)
}
