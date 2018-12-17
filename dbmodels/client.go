package dbmodels

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	// this blank import is needed for the migration script functionality
	_ "github.com/golang-migrate/migrate/source/file"
)

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
func DeleteBookingSystem(hotelID *uuid.UUID, bID uuid.UUID) (int, error) {

	return 0, nil
}
