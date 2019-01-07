package main

import (
	"bookings/dbmodels"
	"bookings/server"
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// Config defines the server configuration options
type Config struct {
	PostgresConfig            dbmodels.Config
	KafkaHosts                []string `envconfig:"kafka_hosts" default:"172.17.42.1:9092"`
	RegionTag                 string   `envconfig:"region_tag" default:"dev"`
	MaxCompressedMessageBytes int      `envconfig:"max_compressed_message_bytes" default:"5000000"`
	KafkaLoggingLevel         string   `envconfig:"kafka_logging_level" default:"warning"`
	KafkaAuditTopic           string   `default:"history"`
}

const (
	// Component is the identifier for this component
	Component = "bookings"
)

func main() {
	var conf Config
	var swaggercapi = flag.Bool("swaggercapi", false, "generate swagger json")
	var swaggerpapi = flag.Bool("swaggerpapi", false, "generate swagger json")
	var swaggersapi = flag.Bool("swaggersapi", false, "generate swagger json")
	var migrate = flag.Bool("migrate", false, "do db migration")
	conf.PostgresConfig = dbmodels.Config{
		DBName:   "bookings",
		Host:     "localhost",
		Password: "bookings",
		Port:     5432,
		User:     "bookings",
	}
	if *swaggercapi {
		api := server.CreateSwaggerCAPI()
		sw, _ := api.RenderJSON()
		fmt.Println(string(sw))
		os.Exit(0)
	}
	if *swaggerpapi {
		api := server.CreateSwaggerPAPI()
		sw, _ := api.RenderJSON()
		fmt.Println(string(sw))
		os.Exit(0)
	}
	if *swaggersapi {
		api := server.CreateSwaggerSAPI()
		sw, _ := api.RenderJSON()
		fmt.Println(string(sw))
		os.Exit(0)
	}
	spew.Println(conf.PostgresConfig)
	database, err := dbmodels.Connect(conf.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to bookings db: %s", err)
	}
	if *migrate {
		dbmodels.Migrate(database, "migrations")
		os.Exit(0)
	}

	log.Info("Starting up Bookings API ...")
	server.RunServer()

	log.Info("Shutting Down")
	os.Exit(0)
}
