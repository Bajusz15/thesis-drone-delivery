package models

import "time"

type Telemetry struct {
	Speed              float64          `json:"speed" db:"speed" bson:"speed"` // meter/sec
	Location           GPS              `json:"location" db:"location" bson:"location"`
	Altitude           float64          `json:"altitude" db:"altitude" bson:"altitude"` // in meters
	Bearing            float64          `json:"bearing" db:"bearing" bson:"bearing"`
	Acceleration       float64          `json:"acceleration" db:"acceleration" bson:"acceleration"`
	BatteryLevel       int              `json:"battery_level" db:"battery_level" bson:"battery_level"` // ez csak százalék, 1-100
	BatteryTemperature int              `json:"battery_temperature" db:"battery_temperature" bson:"battery_temperature"`
	MotorTemperatures  []int            `json:"motor_temperatures" db:"motor_temperatures" bson:"motor_temperatures"`
	Errors             []TelemetryError `json:"errors" db:"errors" bson:"errors"`
	TimeStamp          time.Time        `json:"time_stamp" db:"time_stamp" bson:"time_stamp"`
	DroneID            int              `json:"drone_id" db:"drone_id" bson:"drone_id"`
}

type TelemetryError int

const (
	MotorFailure TelemetryError = iota
	BeaconSignalStrengthLow
	BeaconSignalInterference
	BeaconTemperatureTooHigh
	GPSInteference
	GPSSignalLost
	GPSTemperatureTooHigh
	ProcessorTemperatureTooHigh
	BatteryFailure
	FailedToEjectPackage
	PackageLost
	DestinationDistanceTooFar
)

type GPS struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}
