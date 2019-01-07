package server

import (
	"bookings/dbmodels"
	"bookings/handlers"
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"
)

var itmUUID swagger.Items

func bookingsCAPI() []*swagger.Endpoint {
	itmUUID = swagger.Items{
		Format: "uuid",
		Type:   "string",
	}

	getBookingsCustomer := endpoint.New("GET", "/booking_requests", "Get booking request",
		endpoint.Handler(handlers.GetBookings),
		endpoint.Description("Get all the booking requests per customer"),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Nullable:    true,
				Description: "Page-number to show as first page",
			},
			"per_page": {
				Type:        "integer",
				Nullable:    true,
				Description: "Number of records on a page",
			},
			"data_center": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of data-centers uuids",
			},
			"booking": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of booking uuids",
			},
			"requestor": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of requestors uuids",
			},
			"fromdate": {
				Type:        "string",
				Format:      "date",
				Nullable:    true,
				Description: "the starting date of the listed entries",
			},
			"todate": {
				Type:        "string",
				Format:      "date",
				Nullable:    true,
				Description: "the ending date of the listed entries",
			},
			"states": {
				Type:        "string",
				Enum:        []string{"draft", "cancelled", "approved", "pending", "pending_resp", "rejected", "completed"},
				Nullable:    true,
				Description: "comma separated list of states {'draft', 'cancelled', 'approved', 'pending', 'pending_resp', 'rejected', 'completed'}",
			},
		}),
		endpoint.Response(http.StatusOK, []dbmodels.Booking{}, "Success"),
		endpoint.Tags("Booking Requests CAPI"),
	)
	getBookingCustomer := endpoint.New("GET", "/booking_requests/{id}", "Get booking request",
		endpoint.Handler(handlers.GetBooking),
		endpoint.Description("Get booking request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "booking request id"),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "Success"),
		endpoint.Tags("Booking Requests CAPI"),
	)
	postBookingCustomer := endpoint.New("POST", "/booking_requests", "Create a booking request",
		endpoint.Handler(handlers.PostBooking),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a booking request"),
		endpoint.Body(dbmodels.BookingPost{}, "booking request post body", true),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "SUCCESS"),
		endpoint.Tags("Booking Requests CAPI"),
	)
	patchBookingCustomer := endpoint.New("PATCH", "/booking_requests/{id}", "Update a booking request",
		endpoint.Handler(handlers.PatchBooking),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "booking request id"),
		endpoint.Description("Update a booking request"),
		endpoint.Body(dbmodels.BookingPatch{}, "booking request patch body", true),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "UPDATED"),
		endpoint.Tags("Booking Requests CAPI"),
	)
	return []*swagger.Endpoint{
		getBookingsCustomer,
		getBookingCustomer,
		postBookingCustomer,
		patchBookingCustomer,
	}
}
func bookingsPAPI() []*swagger.Endpoint {
	itmUUID = swagger.Items{
		Format: "uuid",
		Type:   "string",
	}

	getBookingsProvider := endpoint.New("GET", "/provider/booking_requests", "Get booking request",
		endpoint.Handler(handlers.GetBookingsPAPI),
		endpoint.Description("Get all the booking requests per customer"),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Nullable:    true,
				Description: "Page-number to show as first page",
			},
			"per_page": {
				Type:        "integer",
				Nullable:    true,
				Description: "Number of records on a page",
			},
			"data_center": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of data-centers uuids",
			},
			"booking": {
				Type: "array",
				Items: &swagger.Items{
					Format: "uuid",
					Type:   "string"},
				Nullable:    true,
				Description: "comma separated list of booking uuids",
			},
			"customer": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of customer uuids",
			},
			"requestor": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of requestors uuids",
			},
			"fromdate": {
				Type:        "string",
				Format:      "date",
				Nullable:    true,
				Description: "the ending date of the listed entries",
			},
			"todate": {
				Type:        "string",
				Format:      "date",
				Nullable:    true,
				Description: "the starting date of the listed entries",
			},
			"states": {
				Type:        "string",
				Enum:        []string{"draft", "cancelled", "approved", "pending", "pending_resp", "rejected", "completed"},
				Nullable:    true,
				Description: "comma separated list of states {'draft', 'cancelled', 'approved', 'pending', 'pending_resp', 'rejected', 'completed'}",
			},
		}),
		endpoint.Response(http.StatusOK, []dbmodels.Booking{}, "Success"),
		endpoint.Tags("Booking Requests PAPI"),
	)
	getBookingProvider := endpoint.New("GET", "/provider/booking_requests/{id}", "Get booking request",
		endpoint.Handler(handlers.GetBooking),
		endpoint.Description("Get booking request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "booking request id"),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "Success"),
		endpoint.Tags("Booking Requests PAPI"),
	)
	postBookingProvider := endpoint.New("POST", "/provider/booking_requests", "Create a booking request",
		endpoint.Handler(handlers.PostBookingPAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a booking request"),
		endpoint.Body(dbmodels.BookingPost{}, "booking request post body", true),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "SUCCESS"),
		endpoint.Tags("Booking Requests PAPI"),
	)
	patchBookingProvider := endpoint.New("PATCH", "/provider/booking_requests/{id}", "Update a booking request",
		endpoint.Handler(handlers.PatchBookingPAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "booking request id"),
		endpoint.Description("Update a booking request"),
		endpoint.Body(dbmodels.BookingPatch{}, "booking request patch body", true),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "UPDATED"),
		endpoint.Tags("Booking Requests PAPI"),
	)
	return []*swagger.Endpoint{
		getBookingsProvider,
		getBookingProvider,
		postBookingProvider,
		patchBookingProvider,
	}
}
func bookingsSAPI() []*swagger.Endpoint {
	postBookingSystem := endpoint.New("POST", "/system/booking_requests", "Create a booking request",
		endpoint.Handler(handlers.PostBookingSAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a booking request"),
		endpoint.Body(dbmodels.BookingPost{}, "booking request post body", true),
		endpoint.Response(http.StatusOK, dbmodels.Booking{}, "SUCCESS"),
		endpoint.Tags("Booking Requests SAPI"),
	)

	deleteBookingSystem := endpoint.New("DELETE", "/restricted/booking_requests/{id}", "Delete booking request",
		endpoint.Handler(handlers.DeleteBookingSystem),
		endpoint.Description("Delete booking request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "booking request id"),
		endpoint.Response(http.StatusNoContent, "Success", "Successful booking request removal"),
		endpoint.Tags("Booking Requests SAPI"),
	)

	return []*swagger.Endpoint{
		postBookingSystem,
		deleteBookingSystem,
	}
}
