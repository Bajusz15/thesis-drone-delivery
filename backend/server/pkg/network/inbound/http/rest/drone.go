package rest

import (
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/domain/services/drone"
	"drone-delivery/server/pkg/domain/services/telemetry"
	"github.com/labstack/echo/v4"
	"net/http"
)

type FlightData struct {
	Drone       models.Drone       `json:"drone"`
	Telemetries []models.Telemetry `json:"telemetries"`
}

func GetDronesInDelivery(d drone.Service, t telemetry.Service) echo.HandlerFunc {
	return func(context echo.Context) error {
		payload := []FlightData{}
		drones, err := d.GetDronesDelivering()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get drones delivering parcel")
		}

		for _, d := range drones {
			ts, err := t.GetDroneTelemetry(d.ID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get telemetry data of drone")
			}
			fl := FlightData{
				Drone:       d,
				Telemetries: ts,
			}
			payload = append(payload, fl)
		}
		return context.JSON(http.StatusOK, payload)
	}
}
