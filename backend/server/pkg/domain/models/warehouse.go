package models

type Warehouse struct {
	ID       int `json:"id" db:"id" bson:"id"`
	Location GPS `json:"location" db:"location" bson:"location"`
}
