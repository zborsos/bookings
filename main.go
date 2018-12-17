package main

import (
	"bookings/server"
	"bookings/dbmodels"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	// Component is the identifier for this component
	Component = "booking"
)

func main() {
	var conf Config
	var swaggercapi = flag.Bool("swaggercapi", false, "generate swagger json")
	var swaggerpapi = flag.Bool("swaggerpapi", false, "generate swagger json")
	var swaggersapi = flag.Bool("swaggersapi", false, "generate swagger json")
	var migrate = flag.Bool("migrate", false, "do db migration")

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
	
		if *migrate {
			dbmodels.Migrate("migrations")
			os.Exit(0)
		}
	
	log.Info("Starting up Bookings API ...")
	server.RunServer({

	})

	log.Info("Shutting Down")
	os.Exit(0)
}
