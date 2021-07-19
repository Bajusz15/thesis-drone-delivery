package main

import (
	"context"
	"drone-delivery/server/pkg/domain/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func main() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//log.Println("mongodb://"+sc.UserName+":"+sc.PW+"@"+sc.Host+":"+sc.Port+"/"+sc.Database)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://drone-user:drone-pwd@localhost:27017/drone_delivery"))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+sc.Host+":"+sc.Port+"/"+sc.Database))
	if err != nil {
		panic(err)
	}
	//defer client.Disconnect(ctx)
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		panic(err)
	}
	log.Println("You are connected to your database")
	db := client.Database("drone_delivery")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()
	var telemetries []models.Telemetry
	groupStage := bson.D{{"$group", bson.D{{"_id", "$drone_id"}}}}
	sortStage := bson.D{{"$sort", bson.D{{"time_stamp", -1}}}}
	//sortStage := bson.D{bson.M{
	//	"$sort": bson.M{
	//		"time_stamp": -1,
	//	}}}
	cur, err := db.Collection("telemetry").Aggregate(ctx2, mongo.Pipeline{sortStage, groupStage})
	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(ctx2)
	for cur.Next(ctx2) {
		var result models.Telemetry
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(err)
		}
		//telemetries = append(telemetries, result)
	}
	//fmt.Println(telemetries)

	filter := bson.D{}
	//bson.M{
	//	"time_stamp": -1,
	//}}

	opts := options.Distinct().SetMaxTime(2 * time.Second)
	values, err := db.Collection("telemetry").Distinct(context.TODO(), "drone_id", filter, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range values {
		fmt.Println(value)
		var result models.Telemetry
		cur2 := db.Collection("telemetry").FindOne(ctx2, filter)
		err = cur2.Decode(&result)
		telemetries = append(telemetries, result)
	}
	fmt.Println(telemetries)

}
