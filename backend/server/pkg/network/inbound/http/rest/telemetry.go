package rest

import (
	"drone-delivery/server/pkg/domain/models"
	"drone-delivery/server/pkg/domain/services/drone"
	"drone-delivery/server/pkg/domain/services/telemetry"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TelemetryData struct {
	Telemetry models.Telemetry `json:"telemetry"`
}

func SaveTelemetry(d drone.Service, t telemetry.Service) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		var td TelemetryData
		//b, err := ioutil.ReadAll(context.Request().Body)
		//defer context.Request().Body.Close()
		//fmt.Println(string(b))
		//log.Fatal("allj meg")
		//
		//err = json.Unmarshal(b, &td)
		err = context.Bind(&td)

		//log.Println(td)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not save telemetry, ")
		}
		err = t.SaveTelemetry(td.Telemetry)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not save telemetry")
		}
		return context.JSON(http.StatusOK, "succesfully saved telemetry")
	}
}

func GetTelemetry(t telemetry.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		option := c.QueryParams().Get("from")
		switch option {
		case "latest":
			telemetries, err := t.GetLatestTelemetryInECEF()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get telemetry")
			}
			return c.JSON(http.StatusOK, telemetries)
		default:
			telemetries, err := t.GetAllTelemetryInECEF()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get telemetry")
			}
			return c.JSON(http.StatusOK, telemetries)
		}
	}
}
