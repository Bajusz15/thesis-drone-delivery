package mongodb

import (
	"context"
	"drone-delivery/server/pkg/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Storage struct {
	client *mongo.Client
	db     *mongo.Database
}

type StorageConfig struct {
	UserName string
	Database string
	Host     string
	Port     string
	PW       string
}

func NewStorage(sc StorageConfig) (*Storage, error) {
	s := new(Storage)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//log.Println("mongodb://" + sc.UserName + ":" + sc.PW + "@" + sc.Host + ":" + sc.Port + "/" + sc.Database)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+sc.UserName+":"+sc.PW+"@"+sc.Host+":"+sc.Port+"/"+sc.Database))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+sc.Host+":"+sc.Port+"/"+sc.Database))
	if err != nil {
		return nil, err
	}
	s.client = client
	//defer client.Disconnect(ctx)
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		return nil, err
	}
	log.Println("You are connected to your database")
	s.db = s.client.Database(sc.Database)
	_ = s.db.CreateCollection(ctx, "warehouse")
	_ = s.db.CreateCollection(ctx, "telemetry")
	_ = s.db.CreateCollection(ctx, "drone")
	_ = s.db.CreateCollection(ctx, "parcel")
	return s, nil
}

func (s *Storage) CloseConnection() error {
	err := s.client.Disconnect(context.Background())
	return err
}

func (s *Storage) GetWarehouse() (models.Warehouse, error) {
	var wh models.Warehouse
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.db.Collection("warehouse").FindOne(ctx, bson.D{}).Decode(&wh)
	return wh, err
}

func (s *Storage) GetTelemetriesByDrone(droneID int) ([]models.Telemetry, error) {
	var t []models.Telemetry
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := s.db.Collection("telemetry").Find(ctx, bson.D{{"drone_id", droneID}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.Telemetry
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		t = append(t, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Storage) InsertTelemetry(t models.Telemetry) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err != nil {
		return err
	}
	_, err = s.db.Collection("telemetry").InsertOne(ctx, t)
	return err
}

func (s *Storage) GetFreeDrones() ([]models.Drone, error) {
	var drones []models.Drone
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := s.db.Collection("drone").Find(ctx, bson.D{{"state", "free"}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.Drone
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		drones = append(drones, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return drones, nil
}

func (s *Storage) SetDroneState(id int, state string) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err != nil {
		return err
	}
	filter := bson.D{{"id", id}}
	_, err = s.db.Collection("telemetry").UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"state": state,
		},
	})
	return err
}

func (s *Storage) GetParcelsInWarehouse() ([]models.Parcel, error) {
	var parcels []models.Parcel
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := s.db.Collection("parcel").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.Parcel
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return parcels, nil
}

func (s *Storage) GetDronesDelivering() ([]models.Drone, error) {
	var err error
	var drones []models.Drone
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := s.db.Collection("drone").Find(ctx, bson.D{{"state", "in_flight"}})
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.Drone
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		drones = append(drones, result)
	}
	if err != nil {
		return nil, err
	}
	return drones, nil
}

func (s *Storage) GetLatestTelemetryOfDrones() ([]models.Telemetry, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var telemetries []models.Telemetry
	//groupStage := bson.D{{"$group", bson.D{{"_id", "$drone_id"}}}}
	//sortStage := bson.D{{"$sort", bson.D{{"time_stamp", -1}}}}
	//sortStage := bson.D{bson.M{
	//	"$sort": bson.M{
	//		"time_stamp": -1,
	//	}}}
	//cur, err := s.db.Collection("telemetry").Aggregate(ctx, mongo.Pipeline{sortStage, groupStage})
	droneIDs, err := s.db.Collection("telemetry").Distinct(ctx, "drone_id", bson.D{})
	if err != nil {
		return nil, err
	}

	for _, id := range droneIDs {
		var result models.Telemetry
		cur := s.db.Collection("telemetry").FindOne(ctx, bson.D{{"drone_id", id}},
			options.FindOne().SetSort(bson.D{{"time_stamp", -1}}))
		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		telemetries = append(telemetries, result)
	}
	return telemetries, nil
}

func (s *Storage) GetAllTelemetry() ([]models.Telemetry, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var telemetries []models.Telemetry
	cur, err := s.db.Collection("telemetry").Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"time_stamp", -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result models.Telemetry
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		telemetries = append(telemetries, result)
	}
	if err != nil {
		return nil, err
	}
	return telemetries, nil
}

func (s *Storage) ReInitializeDeliveryData(drones []models.Drone, parcels []models.Parcel) error {
	var err error
	_ = s.db.Collection("telemetry").Drop(context.TODO())
	_ = s.db.Collection("drone").Drop(context.TODO())
	_ = s.db.Collection("parcel").Drop(context.TODO())
	err = s.db.CreateCollection(context.TODO(), "warehouse")
	err = s.db.CreateCollection(context.TODO(), "telemetry")
	err = s.db.CreateCollection(context.TODO(), "drone")
	err = s.db.CreateCollection(context.TODO(), "parcel")
	//db.warehouse.insertOne(     { id: 1,  location: { latitude: 48.080922, longitude: 20.766208} } )
	for _, d := range drones {
		_, err = s.db.Collection("drone").InsertOne(context.TODO(), d)
		if err != nil {
			return err
		}
	}

	for _, p := range parcels {
		_, err = s.db.Collection("parcel").InsertOne(context.TODO(), p)
		if err != nil {
			return err
		}
	}
	return nil
}
