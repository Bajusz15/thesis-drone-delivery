package routing

import "github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"

type Service interface {
	UpdateRoute() error
	CalculateDroneDistanceAndDirectionFromDestination(currentLat, currentLon, destinationLat, destinationLon float64) (distance, bearing float64)
	CalculateDroneNextCoordinates(lat, lon, dist, bearing float64) (nextLat, nextLon float64)
}

type service struct {
	geometry ellipsoid.Ellipsoid
}

func (s *service) UpdateRoute() error {
	panic("implement me")
}

func NewService(geo ellipsoid.Ellipsoid) *service {
	return &service{geo}
}

func (s *service) CalculateDroneDistanceAndDirectionFromDestination(currentLat, currentLon, destinationLat, destinationLon float64) (distance, bearing float64) {
	distance, bearing = s.geometry.To(currentLat, currentLon, destinationLat, destinationLon)
	return distance, bearing
}

func (s *service) CalculateDroneNextCoordinates(lat, lon, dist, bearing float64) (nextLat, nextLon float64) {
	// Calculate where you are when going from  lat/long with direction and distance
	nextLat, nextLon = s.geometry.At(lat, lon, dist, bearing)
	return nextLat, nextLon
}
