package dbmodels

import (
	"time"
	"workflow"

	uuid "github.com/satori/go.uuid"
)

// BookingsResponse ...
type BookingsResponse struct {
	NumResults int       `json:"num_results"`
	Objects    []Booking `json:"objects"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
}

// Booking struct
type Booking struct {
	ID                      uuid.UUID             `db:"id" json:"id"`
	RoomID                  uuid.UUID             `db:"room_id" json:"room_id"`
	CustomerID              uuid.UUID             `db:"customer_id" json:"customer_id"`
	RequestorID             uuid.UUID             `db:"requestor_id" json:"requestor_id"`
	RequestedAt             time.Time             `db:"requested_at" json:"requested_at"`
	StartTime               time.Time             `db:"start_time" json:"start_time"`
	EndTime                 time.Time             `db:"end_time" json:"end_time"`
	State                   string                `db:"state" json:"state"`
	StateInfo               *string               `db:"state_information" json:"state_information"`
	BucketName              *string               `db:"file_name" json:"bucket_name"`
	Description             *string               `db:"description" json:"description"`
	Reference               *string               `db:"reference" json:"reference"`
	Transitions             *workflow.Transitions `db:"-" json:"transitions"`
	BookingRequestEmail     *string               `db:"booking_request_email" json:"booking_request_email,omitempty"`
	BookingRequestFromEmail *string               `db:"booking_request_from_email" json:"booking_request_from_email,omitempty"`
}

// Room struct
type Room struct {
	ID                  uuid.UUID `db:"id" json:"id"`
	Name                string    `db:"name" json:"name"`
	HotelID             uuid.UUID `db:"hotel_id" json:"hotel_id"`
	Type                string    `db:"type" json:"type"`
	ReservationMaxTime  *string   `db:"reservation_max_time" json:"reservation_max_time_hours"`
	AvailableFrom       string    `db:"available_from" json:"available_from"`
	AvailableTo         string    `db:"available_to" json:"available_to"`
	ReservationLeadTime *string   `db:"reservation_lead_time" json:"reservation_lead_time_days"`
	IsShared            bool      `db:"is_shared" json:"is_shared"`
	SharedNrPerson      *int16    `db:"shared_nr_person" json:"shared_nr_person"`
	Description         *string   `db:"description" json:"description"`
}

// RoomAvailability shows periods, when the room is booked
type RoomAvailability struct {
	ID        *uuid.UUID `db:"id" json:"request_id"`
	StartTime *time.Time `db:"start_time" json:"start_time"`
	EndTime   *time.Time `db:"end_time" json:"end_time"`
}

// BookingPost ...
type BookingPost struct {
	RoomID                  uuid.UUID             `db:"room_id" json:"room_id"`
	CustomerID              uuid.UUID             `db:"customer_id" json:"customer_id"`
	RequestorID             uuid.UUID             `db:"requestor_id" json:"requestor_id"`
	RequestedAt             time.Time             `db:"requested_at" json:"requested_at"`
	StartTime               time.Time             `db:"start_time" json:"start_time"`
	EndTime                 time.Time             `db:"end_time" json:"end_time"`
	State                   string                `db:"state" json:"state"`
	StateInfo               *string               `db:"state_information" json:"state_information"`
	BucketName              *string               `db:"file_name" json:"bucket_name"`
	Description             *string               `db:"description" json:"description"`
	Reference               *string               `db:"reference" json:"reference"`
	Transitions             *workflow.Transitions `db:"-" json:"transitions"`
	BookingRequestEmail     *string               `db:"booking_request_email" json:"booking_request_email,omitempty"`
	BookingRequestFromEmail *string               `db:"booking_request_from_email" json:"booking_request_from_email,omitempty"`
}

// BookingPatch ...
type BookingPatch struct {
	RoomID                  *uuid.UUID            `db:"room_id" json:"room_id"`
	RequestorID             *uuid.UUID            `db:"requestor_id" json:"requestor_id"`
	RequestedAt             *time.Time            `db:"requested_at" json:"requested_at"`
	StartTime               *time.Time            `db:"start_time" json:"start_time"`
	EndTime                 *time.Time            `db:"end_time" json:"end_time"`
	State                   *string               `db:"state" json:"state"`
	StateInfo               *string               `db:"state_information" json:"state_information"`
	BucketName              *string               `db:"file_name" json:"bucket_name"`
	Description             *string               `db:"description" json:"description"`
	Reference               *string               `db:"reference" json:"reference"`
	Transitions             *workflow.Transitions `db:"-" json:"transitions"`
	BookingRequestEmail     *string               `db:"booking_request_email" json:"booking_request_email,omitempty"`
	BookingRequestFromEmail *string               `db:"booking_request_from_email" json:"booking_request_from_email,omitempty"`
}
