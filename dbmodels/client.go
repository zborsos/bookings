package dbmodels

import (
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	// this blank import is needed for the migration script functionality
	_ "github.com/golang-migrate/migrate/source/file"

	log "github.com/sirupsen/logrus"
)

// Config defines the config information to be passed to Connect method
type Config struct {
	User     string `default:"bookings"`
	DBName   string `default:"bookings"`
	Password string `vaultconfig:"secret/postgres/bookings"`
	Host     string `envconfig:"postgres_host" default:"localhost"`
	Port     int    `envconfig:"postgres_port" default:"5432"`
}

// Connect to PostGres on Bespin returns db client
func Connect(conf Config) (rawdb *sqlx.DB, err error) {
	cn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.User,
		conf.Password)

	rawdb, err = sqlx.Connect("postgres", cn)
	if err != nil {
		log.Errorf("Failed to connect to postgres database: %s", err)
		return
	}
	return
}

// GetBooking ...
func GetBooking(id uuid.UUID, actor string) (*Booking, int, error) {

	return nil, http.StatusOK, nil
}

// GetBookings ...
func GetBookings(actor string) ([]Booking, int, error) {
	return nil, 0, nil
}

// PostBooking ...
func PostBooking(body *BookingPost) (*Booking, int, error) {
	return nil, http.StatusOK, nil
}

// PatchBooking ...
func PatchBooking(body *BookingPatch, id uuid.UUID) (*Booking, int, error) {

	return nil, http.StatusOK, nil
}

// DeleteBookingSystem ...
func DeleteBookingSystem(hotelID *uuid.UUID, bID uuid.UUID) error {

	return nil
}

// Migrate does db migration up to the latest version
func Migrate(db *sqlx.DB, path string) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Errorf("Failed to to get postgres driver: %s", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", driver)
	if err != nil {
		log.Errorf("Failed to create migration client: %s", err)
		return
	}

	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Info("DB is up to date")
		} else {
			log.Errorf("Failed to migrate db: %s", err)
		}
	}
}
