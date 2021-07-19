package postgres

import (
	"database/sql"
	"drone-delivery/server/pkg/domain/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sqlx.DB
}

type StorageConfig struct {
	UserName string
	Database string
	Host     string
	Port     string
	SSSLMode string
	PW       string
}

func NewStorage(sc StorageConfig) (*Storage, error) {
	var err error
	s := new(Storage)
	connStr := "user=" + sc.UserName + " dbname=" + sc.Database + " password=" + sc.PW + " host=" + sc.Host + " sslmode=disable TimeZone=Europe/Budapest"
	s.db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = s.db.Ping(); err != nil {
		return nil, err
	}
	log.Println("You are connected to your database")
	return s, nil
}

func (s *Storage) CloseConnection() error {
	err := s.db.Close()
	return err
}

func NewMockStorage(db *sql.DB) (*Storage, error) {
	s := new(Storage)
	s.db = sqlx.NewDb(db, "postgres")
	return s, nil
}
func (s *Storage) GetWarehouse() (models.Warehouse, error) {
	var wh models.Warehouse
	err := s.db.Get(&wh, `SELECT id, location_latitude "location.latitude", location_longitude "location.longitude" FROM warehouse LIMIT 1`)
	return wh, err
}

func (s *Storage) GetTelemetriesByDrone(droneID int) ([]models.Telemetry, error) {
	var t []models.Telemetry
	err := s.db.Select(&t, `SELECT id, drone_id, speed, longitude "location.longitude",latitude "location.latitude", altitude, bearing, acceleration, battery_level, battery_temperature, motor_temperatures,time_stamp 
									FROM telemetry 
									WHERE drone_id=$1`, droneID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Storage) InsertTelemetry(t models.Telemetry) error {
	motorTemps := pq.Array(t.MotorTemperatures)
	_, err := s.db.Exec(`INSERT INTO telemetry 
    							(drone_id, speed, latitude, longitude, altitude, bearing, acceleration, battery_level,battery_temperature, motor_temperatures, time_stamp)
                       			VALUES ($1, $2, $3,$4,$5,$6,$7, $8, $9, $10, $11) `,
		t.DroneID, t.Speed, t.Location.Latitude, t.Location.Longitude, t.Altitude, t.Bearing, t.Acceleration, t.BatteryLevel,
		t.BatteryTemperature, motorTemps, t.TimeStamp)
	return err
}

func (s *Storage) GetFreeDrones() ([]models.Drone, error) {
	var drones []models.Drone
	var err error
	//tx, err := s.db.Beginx()
	//if err != nil {
	//	return nil, err
	//}
	//this is needed so there is no dirty read, repeatable read or phantom read
	//_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	//if err != nil {
	//	return nil, err
	//}
	err = s.db.Select(&drones, `SELECT id "drone_id", consumption, weight FROM drone WHERE state='free'`)
	if err != nil {
		//tx.Rollback()
		return nil, err
	}
	for _, d := range drones {
		log.Println(d)
	}
	//tx.Commit()
	return drones, nil
}

func (s *Storage) SetDroneState(id int, state string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	//this is needed so there is no dirty read, repeatable read or phantom read
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE drone SET state=$1 WHERE id=$2 ", state, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *Storage) GetParcelsInWarehouse() ([]models.Parcel, error) {
	var parcels []models.Parcel
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	//this is needed so there is no dirty read, repeatable read or phantom read
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		return nil, err
	}
	err = tx.Select(&parcels, `SELECT id, name, tracking_id, weight, drop_off_latitude "drop_off_site.latitude", drop_off_longitude "drop_off_site.longitude" FROM parcel`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return parcels, nil
}

func (s *Storage) GetDronesDelivering() ([]models.Drone, error) {
	var err error
	var drones []models.Drone
	err = s.db.Select(&drones, `SELECT id, state, weight, consumption FROM drone WHERE state=$1`, models.DroneInFlight)
	if err != nil {
		return nil, err
	}
	return drones, nil
}

func (s *Storage) GetAllTelemetry() ([]models.Telemetry, error) {
	var err error
	var telemetries []models.Telemetry
	err = s.db.Select(&telemetries, `SELECT  drone_id, speed, latitude "location.latitude", longitude "location.longitude", altitude,
       bearing, acceleration, battery_level, battery_temperature, time_stamp 
		FROM telemetry 
		ORDER BY time_stamp`)
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}

func (s *Storage) GetLatestTelemetryOfDrones() ([]models.Telemetry, error) {
	var err error
	var telemetries []models.Telemetry
	//var droneIDs []int
	//for _, d := range drones {
	//	droneIDs = append(droneIDs, d.ID)
	//}

	err = s.db.Select(&telemetries, `SELECT  DISTINCT ON (drone_id) drone_id, speed, latitude "location.latitude", longitude "location.longitude", altitude,
       bearing, acceleration, battery_level, battery_temperature, time_stamp 
		FROM telemetry  
		ORDER BY drone_id, time_stamp DESC`)
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}

func (s *Storage) ReInitializeDeliveryData(drones []models.Drone, parcels []models.Parcel) error {
	var err error
	_, err = s.db.Exec(`TRUNCATE drone RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`TRUNCATE parcel RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`TRUNCATE telemetry RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}

	for _, d := range drones {
		_, err = s.db.Exec(`INSERT INTO drone (id, state, weight, consumption) 
			VALUES ($1, $2, $3, $4)`, d.ID, d.State, d.Weight, d.Consumption)
		if err != nil {
			return err
		}
	}

	for _, p := range parcels {
		_, err = s.db.Exec(`INSERT INTO parcel (id, name, weight, drop_off_latitude, drop_off_longitude ) 
				VALUES ($1, $2, $3, $4, $5)`, p.ID, p.Name, p.Weight, p.DropOffSite.Latitude, p.DropOffSite.Longitude)
		if err != nil {
			return err
		}
	}
	return nil
}
