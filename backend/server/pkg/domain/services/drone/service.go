package drone

import (
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/domain/services/routing"
	"errors"
	gokitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"log"
	"math/rand"
	"time"
)

type Service interface {
	DeliverParcels() error
	ProvisionDrone(wh models.Warehouse, d models.Drone) error
	GetFreeDrones() ([]models.Drone, error)
	GetDronesDelivering() ([]models.Drone, error)
	ChangeService(r Repository)
	ReinitializeDatabase(repos ...Repository) error
}

type Repository interface {
	GetFreeDrones() ([]models.Drone, error)
	GetParcelsInWarehouse() ([]models.Parcel, error)
	GetWarehouse() (models.Warehouse, error)
	GetDronesDelivering() ([]models.Drone, error)
	SetDroneState(droneID int, state string) error
	ReInitializeDeliveryData(drones []models.Drone, parcels []models.Parcel) error
}

type OutboundAdapter interface {
	FetchProvisionDroneEndpoint(wh models.Warehouse, d models.Drone) (success bool, err error)
}

type service struct {
	repo           Repository
	adapter        OutboundAdapter
	logger         gokitlog.Logger
	routingService routing.Service
}

func (s *service) GetDronesDelivering() ([]models.Drone, error) {
	panic("implement me")
}

func NewService(r Repository, ea OutboundAdapter, l gokitlog.Logger, rs routing.Service) *service {
	return &service{r, ea, l, rs}
}

func (s *service) ChangeService(r Repository) {
	s.repo = r
}

func (s *service) DeliverParcels() error {
	wh, err := s.repo.GetWarehouse()
	if err != nil {
		s.logger.Log("desc", "could not get warehouse from repository", "err", err)
		return err
	}

	parcels, err := s.repo.GetParcelsInWarehouse()
	if err != nil {
		s.logger.Log("err", err)
	}
	freeDrones, err := s.GetFreeDrones()
	if err != nil {
		s.logger.Log("desc", "could not get free drones from repository", "err", err)
		return err
	}
	log.Println(parcels)
	log.Println(freeDrones)
	drones, err := s.routingService.OptimizeRoutes(wh, freeDrones, parcels)
	if err != nil {
		return errors.New("failed to set routes for drones, aborting delivery")
	}

	for _, d := range drones {
		err = s.ProvisionDrone(wh, d)
		if err != nil {
			s.logger.Log("desc", "could not provision drone")
			continue
		}
		err = s.repo.SetDroneState(d.ID, "in-flight")
		if err != nil {
			s.logger.Log("err", err, "desc", "failed to set drone state")
			continue
		}

	}

	return err
}

func (s *service) ProvisionDrone(wh models.Warehouse, d models.Drone) error {
	logger := gokitlog.With(s.logger, "method", "ProvisionDrone")
	success, err := s.adapter.FetchProvisionDroneEndpoint(wh, d)
	if err != nil {
		level.Warn(logger).Log(
			"description", "could not provision drone, outbound adapter returned an error",
			"err", err,
		)
	}
	if !success {
		return errors.New("could not start drone")
	}

	return nil
}

func (s *service) GetFreeDrones() ([]models.Drone, error) {
	logger := gokitlog.With(s.logger, "method", "GetFreeDrones")
	drones, err := s.repo.GetFreeDrones()
	if err != nil {
		level.Error(logger).Log(
			"desc", "could not get drones, repository returned an error",
			"err", err,
		)
	}
	return drones, nil
}

func (s *service) ReinitializeDatabase(repos ...Repository) error {
	var err error
	logger := gokitlog.With(s.logger, "method", "ReinitializeDatabase")
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30
	deliveries := rand.Intn(max-min+1) + min
	//generate latitudes and longitudes for the parcels
	latitudes := randFloats(48.05, 48.10, deliveries)
	longitudes := randFloats(20.73, 20.78, deliveries)
	parcelWeights := randFloats(0.2, 2, deliveries)
	consumptions := randFloats(300, 800, deliveries)
	droneWeights := randFloats(3, 15, deliveries)

	var drones = make([]models.Drone, deliveries)
	var parcels = make([]models.Parcel, deliveries)
	for i := 0; i < deliveries; i++ {
		drones[i] = models.Drone{
			ID:          i + 1,
			Consumption: consumptions[i],
			Weight:      droneWeights[i],
			State:       models.DroneFree,
		}

		parcels[i] = models.Parcel{
			ID:     i + 1,
			Name:   "egy csomag",
			Weight: parcelWeights[i],
			DropOffSite: models.GPS{
				Latitude:  latitudes[i],
				Longitude: longitudes[i],
			},
		}
	}
	for _, r := range repos {
		err = r.ReInitializeDeliveryData(drones, parcels)
		if err != nil {
			level.Error(logger).Log(
				"desc", "could not reinitialize db",
				"err", err,
			)
			return err
		}
	}

	return nil
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}
