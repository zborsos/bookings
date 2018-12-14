package server

import (
	"booking/handlers"
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

	getBookingsCustomer := endpoint.New("GET", "/facility_requests", "Get facility request",
		endpoint.Handler(handlers.GetBookings),
		endpoint.Description("Get all the facility requests per customer"),
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
			"facility": {
				Type:        "array",
				Items:       &itmUUID,
				Nullable:    true,
				Description: "comma separated list of facilities uuids",
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
		endpoint.Response(http.StatusOK, []db.Booking{}, "Success"),
		endpoint.Tags("Facility Requests CAPI"),
	)
	getBookingCustomer := endpoint.New("GET", "/facility_requests/{id}", "Get facility request",
		endpoint.Handler(handlers.GetBooking),
		endpoint.Description("Get facility request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "facility request id"),
		endpoint.Response(http.StatusOK, db.Booking{}, "Success"),
		endpoint.Tags("Facility Requests CAPI"),
	)
	postBookingCustomer := endpoint.New("POST", "/facility_requests", "Create a facility request",
		endpoint.Handler(handlers.PostBooking),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a facility request"),
		endpoint.Body(db.BookingPost{}, "facility request post body", true),
		endpoint.Response(http.StatusOK, db.Facility{}, "SUCCESS"),
		endpoint.Tags("Facility Requests CAPI"),
	)
	patchBookingCustomer := endpoint.New("PATCH", "/facility_requests/{id}", "Update a facility request",
		endpoint.Handler(handlers.PatchBooking),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "facility request id"),
		endpoint.Description("Update a facility request"),
		endpoint.Body(db.BookingPatch{}, "facility request patch body", true),
		endpoint.Response(http.StatusOK, db.Booking{}, "UPDATED"),
		endpoint.Tags("Facility Requests CAPI"),
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

	getBookingsProvider := endpoint.New("GET", "/provider/facility_requests", "Get facility request",
		endpoint.Handler(handlers.GetBookingsPAPI),
		endpoint.Description("Get all the facility requests per customer"),
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
			"facility": {
				Type: "array",
				Items: &swagger.Items{
					Format: "uuid",
					Type:   "string"},
				Nullable:    true,
				Description: "comma separated list of facilities uuids",
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
		endpoint.Response(http.StatusOK, []db.Booking{}, "Success"),
		endpoint.Tags("Facility Requests PAPI"),
	)
	getBookingProvider := endpoint.New("GET", "/provider/facility_requests/{id}", "Get facility request",
		endpoint.Handler(handlers.GetBookingPAPI),
		endpoint.Description("Get facility request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "facility request id"),
		endpoint.Response(http.StatusOK, db.Booking{}, "Success"),
		endpoint.Tags("Facility Requests PAPI"),
	)
	postBookingProvider := endpoint.New("POST", "/provider/facility_requests", "Create a facility request",
		endpoint.Handler(handlers.PostBookingPAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a facility request"),
		endpoint.Body(db.BookingPost{}, "facility request post body", true),
		endpoint.Response(http.StatusOK, db.Facility{}, "SUCCESS"),
		endpoint.Tags("Facility Requests PAPI"),
	)
	patchBookingProvider := endpoint.New("PATCH", "/provider/facility_requests/{id}", "Update a facility request",
		endpoint.Handler(handlers.PatchBookingPAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "facility request id"),
		endpoint.Description("Update a facility request"),
		endpoint.Body(db.BookingPatch{}, "facility request patch body", true),
		endpoint.Response(http.StatusOK, db.Booking{}, "UPDATED"),
		endpoint.Tags("Facility Requests PAPI"),
	)
	return []*swagger.Endpoint{
		getBookingsProvider,
		getBookingProvider,
		postBookingProvider,
		patchBookingProvider,
	}
}
func bookingsSAPI() []*swagger.Endpoint {
	postBookingSystem := endpoint.New("POST", "/system/facility_requests", "Create a facility request",
		endpoint.Handler(handlers.PostBookingSAPI),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Description("Create a facility request"),
		endpoint.Body(db.BookingPost{}, "facility request post body", true),
		endpoint.Response(http.StatusOK, db.Facility{}, "SUCCESS"),
		endpoint.Tags("Facility Requests SAPI"),
	)

	deleteBookingSystem := endpoint.New("DELETE", "/restricted/facility_requests/{id}", "Delete facility request",
		endpoint.Handler(handlers.DeleteBookingSystem),
		endpoint.Description("Delete facility request by its ID"),
		endpoint.Query("dc_id", "string", "uuid", "data_center id", false),
		endpoint.Path("id", "string", "uuid", "facility request id"),
		endpoint.Response(http.StatusNoContent, "Success", "Successful facility request removal"),
		endpoint.Tags("Facility Requests SAPI"),
	)

	return []*swagger.Endpoint{
		postBookingSystem,
		deleteBookingSystem,
	}
}
