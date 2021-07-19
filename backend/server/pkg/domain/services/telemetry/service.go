package telemetry

import (
	"drone-delivery/server/pkg/domain/models"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	"github.com/go-kit/kit/log"
)

type Service interface {
	SaveTelemetry(t models.Telemetry) error
	GetDroneTelemetry(droneID int) ([]models.Telemetry, error)
	ChangeService(r Repository)
	GetAllTelemetryInECEF() (map[int][]models.Telemetry, error)
	GetLatestTelemetryInECEF() (map[int][]models.Telemetry, error)
}

type Repository interface {
	InsertTelemetry(t models.Telemetry) error
	GetTelemetriesByDrone(droneID int) ([]models.Telemetry, error)
	GetAllTelemetry() ([]models.Telemetry, error)
	GetLatestTelemetryOfDrones() ([]models.Telemetry, error)
	GetDronesDelivering() ([]models.Drone, error)
}

type service struct {
	repo   Repository
	logger log.Logger
	geo    ellipsoid.Ellipsoid
}

func NewService(r Repository, l log.Logger, g ellipsoid.Ellipsoid) *service {
	return &service{repo: r, logger: l, geo: g}
}

func (s *service) ChangeService(r Repository) {
	s.repo = r
}

func (s *service) SaveTelemetry(t models.Telemetry) error {
	var err error
	//level.Info(s.logger).Log("desc", "saving telemetry", "telemetry", t)
	err = s.repo.InsertTelemetry(t)
	if err != nil {
		s.logger.Log("err", err, "desc", "failed to save drone telemetry")
		return err
	}
	return nil
}

func (s *service) GetDroneTelemetry(droneID int) ([]models.Telemetry, error) {
	telemetries, err := s.repo.GetTelemetriesByDrone(droneID)
	if err != nil {
		s.logger.Log("err", err, "desc", "failed to get drone telemetry")
		return nil, err
	}
	return telemetries, nil
}

func (s *service) GetLatestTelemetryInECEF() (map[int][]models.Telemetry, error) {
	//drones, err := s.repo.GetDronesDelivering()
	//if err != nil {
	//	s.logger.Log("err", err, "desc", "failed to get drones that are delivering")
	//	return nil, err
	//}
	telemetries, err := s.repo.GetLatestTelemetryOfDrones()
	if err != nil {
		s.logger.Log("err", err, "desc", "failed to get all telemetry")
		return nil, err
	}
	telemetryMap := map[int][]models.Telemetry{}
	for _, t := range telemetries {
		t.Location.Latitude, t.Location.Longitude, t.Altitude = s.geo.ToECEF(t.Location.Latitude, t.Location.Longitude, t.Altitude)
		telemetryMap[t.DroneID] = append(telemetryMap[t.DroneID], t)
	}

	return telemetryMap, nil
}

func (s *service) GetAllTelemetryInECEF() (map[int][]models.Telemetry, error) {
	//drones, err := s.repo.GetDronesDelivering()
	//if err != nil {
	//	s.logger.Log("err", err, "desc", "failed to get drones that are delivering")
	//	return nil, err
	//}

	telemetries, err := s.repo.GetAllTelemetry()
	if err != nil {
		s.logger.Log("err", err, "desc", "failed to get all telemetry")
		return nil, err
	}
	telemetryMap := map[int][]models.Telemetry{}
	for _, t := range telemetries {
		t.Location.Latitude, t.Location.Longitude, t.Altitude = s.geo.ToECEF(t.Location.Latitude, t.Location.Longitude, t.Altitude)
		telemetryMap[t.DroneID] = append(telemetryMap[t.DroneID], t)
	}

	return telemetryMap, nil
}
