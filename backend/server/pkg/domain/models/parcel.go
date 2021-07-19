package models

type Parcel struct {
	ID            int             `json:"id" db:"id" bson:"id"`
	TrackingID    string          `json:"tracking_id" db:"tracking_id" bson:"tracking_id"`
	Name          string          `json:"name" db:"name" bson:"name"`
	Weight        float64         `json:"weight" db:"weight"  bson:"weight"`
	Location      GPS             `json:"location" bson:"location"`
	FromAddress   ShippingAddress `json:"from_address"` //ez lehet nem is kell
	ToAddress     ShippingAddress `json:"to_address"`
	DropOffSite   GPS             `bson:"drop_off_site" json:"drop_off_site" db:"drop_off_site"`
	AssignedDrone int             `json:"assigned_drone" db:"assigned_drone" db:"assigned_drone"`
}

type ShippingAddress struct {
	Name    string  `json:"name" validate:"required"`
	Country string  `json:"country" validate:"required"`
	Region  *string `json:"region"`
	City    string  `json:"city" validate:"required"`
	Zip     string  `json:"zip" validate:"required"`
	Street  string  `json:"street" validate:"required"`
	Street2 string  `json:"street_2"`
	Street3 string  `json:"street_3"`
}
