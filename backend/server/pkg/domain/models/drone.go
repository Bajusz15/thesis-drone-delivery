package models

type Drone struct {
	ID           int           `json:"id" db:"drone_id" bson:"id"`
	Telemetry    Telemetry     `json:"telemetry" bson:"telemetry"`
	Parcel       Parcel        `json:"parcel"`
	Destinations []Destination `json:"destinations"`
	Consumption  float64       `json:"consumption" db:"consumption" bson:"consumption"` // electricity used for the drone to travel 1 km with X parcel weight with speed of 10 m/s
	Weight       float64       `json:"weight" db:"weight" bson:"weight"`
	State        DroneState    `db:"state" bson:"state"`
}

type DroneState string

const (
	DroneFree     DroneState = "free"
	DroneInFlight DroneState = "in-flight"
)

func (d Drone) GetConsumption(p Parcel) float64 {
	//most csak 1 csomagot szallit 1 dron, de kesobb akar tobbet is lehet
	if p.Weight == 0 {
		return d.Consumption
	}
	//TODO: ennek biztos van jobb képlete, utána kell nézni
	return p.Weight/d.Weight*200 + d.Consumption
}
