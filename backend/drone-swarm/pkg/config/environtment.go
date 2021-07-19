package config

import (
	"flag"
	"os"
)

var (
	ServerHTTPDomain string
	ServerHTTPPort   string
	ServerGRPCDomain string
	ServerGRPCPort   string
)

var PostgresConfig struct {
	UserName string
	Database string
	Host     string
	Port     string
	SSSLMode string
	PW       string
}

func SetConfig() {
	flag.StringVar(&ServerHTTPDomain, "server domain for http", os.Getenv("SERVER_DOMAIN"), "domain of server, with protocol")
	flag.StringVar(&ServerHTTPPort, "server port for http", os.Getenv("SERVER_PORT"), "the port the server is listening on")
	flag.StringVar(&ServerGRPCDomain, "server domain for grpc", os.Getenv("SERVER_GRPC_DOMAIN"), "domain of server, with protocol")
	flag.StringVar(&ServerGRPCPort, "server port for grpc", os.Getenv("SERVER_GRPC_PORT"), "the port the server is listening on")

}
