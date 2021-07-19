package config

import (
	"flag"
	"os"
)

var DroneSwarmURL string

var PostgresConfig struct {
	UserName string
	Database string
	Host     string
	Port     string
	SSSLMode string
	PW       string
}

var MongoConfig struct {
	UserName string
	Database string
	Host     string
	Port     string
	PW       string
}

func SetConfig() {
	flag.StringVar(&DroneSwarmURL, "drone swarm domain url", os.Getenv("DRONE_SWARM_URL"),
		"An url for the drone swarm application, with protocol")

	PostgresConfig.UserName = os.Getenv("PGUSER")
	PostgresConfig.Database = os.Getenv("PGDATABASE")
	PostgresConfig.Host = os.Getenv("PGHOST")
	PostgresConfig.Port = os.Getenv("PGPORT")
	PostgresConfig.SSSLMode = os.Getenv("PGSSL")
	PostgresConfig.PW = os.Getenv("PGPASSWORD")

	MongoConfig.UserName = os.Getenv("MONGO_USER")
	MongoConfig.Database = os.Getenv("MONGO_DB")
	MongoConfig.Host = os.Getenv("MONGO_HOST")
	MongoConfig.Port = os.Getenv("MONGO_PORT")
	MongoConfig.PW = os.Getenv("MONGO_PWD")
}
