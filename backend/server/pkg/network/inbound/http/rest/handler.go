package rest

import (
	"drone-delivery/server/pkg/domain/services/drone"
	"drone-delivery/server/pkg/domain/services/telemetry"
	"drone-delivery/server/pkg/storage/mongodb"
	"drone-delivery/server/pkg/storage/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Handler(d drone.Service, t telemetry.Service, p *postgres.Storage, m *mongodb.Storage) http.Handler {
	router := echo.New()
	router.Use(middleware.CORS())
	router.POST("/api/delivery", func(c echo.Context) error {
		err := d.DeliverParcels()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to start delivery")
		}
		return c.JSON(http.StatusOK, "delivery started")
	})

	router.GET("/api/delivery/telemetry", GetTelemetry(t))
	router.POST("/api/delivery/telemetry", SaveTelemetry(d, t))
	router.GET("/api/delivery/drones", GetDronesInDelivery(d, t))
	router.PUT("/configure/database/:name", func(c echo.Context) error {
		switch c.Param("name") {
		case "mongodb":
			t.ChangeService(m)
			d.ChangeService(m)
		case "postgres":
			t.ChangeService(p)
			d.ChangeService(p)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "no such database supported")
		}
		return c.JSON(http.StatusOK, "configuration complete")
	})

	router.POST("/api/delivery/reinitialize", func(c echo.Context) error {
		err := d.ReinitializeDatabase(p, m)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to start reinitialize db")
		}
		return c.JSON(http.StatusOK, "initialization finished, db is filled up with drones and parcels")
	})
	return router
}
