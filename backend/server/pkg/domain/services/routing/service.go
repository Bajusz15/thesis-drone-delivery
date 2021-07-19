package routing

import (
	"drone-delivery/server/pkg/domain/models"
	goKitLog "github.com/go-kit/kit/log"
	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
	"math"
)

type Service interface {
	CalculateDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64
	OptimizeRoutes(warehouse models.Warehouse, drones []models.Drone, parcels []models.Parcel) ([]models.Drone, error)
}

type service struct {
	logger goKitLog.Logger
}

func NewService(l goKitLog.Logger) *service {
	return &service{l}
}

func (s *service) OptimizeRoutes(warehouse models.Warehouse, drones []models.Drone, parcels []models.Parcel) ([]models.Drone, error) {
	size := len(drones)
	size2 := len(parcels)
	var b = make([][]int, size)
	for i, d := range drones {
		b[i] = make([]int, size2)
		for j, p := range parcels {
			cost := d.GetConsumption(p) * s.CalculateDistance(
				warehouse.Location.Latitude,
				warehouse.Location.Longitude,
				p.DropOffSite.Latitude,
				p.DropOffSite.Longitude,
				"K",
			)
			b[i][j] = int(math.Round(cost * 100))
		}
	}

	solution, err := hungarianAlgorithm.Solve(b)
	if err != nil {
		s.logger.Log("err", err, "desc", "failed to optimize routes")
		return nil, err
	}

	warehouseDestination := models.Destination{
		Coordinates:          warehouse.Location,
		ParcelDestination:    false,
		WarehouseDestination: true,
	}

	for i := range drones {
		drones[i].Parcel = parcels[solution[i]]
		parcelDestination := models.Destination{
			Coordinates:          parcels[solution[i]].DropOffSite,
			ParcelDestination:    true,
			WarehouseDestination: false,
		}
		//here we could append more destinations if there is something blocking the straight route
		drones[i].Destinations = append(drones[i].Destinations, parcelDestination)
		//last destination is the warehouse
		drones[i].Destinations = append(drones[i].Destinations, warehouseDestination)
	}
	return drones, nil
}

func (s *service) CalculateDistance(lat1, lng1, lat2, lng2 float64, unit ...string) float64 {
	const PI = float64(math.Pi) //3.141592653589793

	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180

	theta := lng1 - lng2
	radtheta := PI * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		} else if unit[0] == "METER" {
			dist = dist * 1.609344 * 1000
		}
	}

	return dist
}
