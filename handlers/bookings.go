package controllers

import (
	"errors"
	"facility_request/db"
	"fmt"
	"net/http"
	"strings"
	"time"

	middlewares "git.ntteo.net/go-libs.git/gin-middlewares"
	"git.ntteo.net/go-libs.git/workflow"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/yaml"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var defaultLang = "en-GB"
var myI18n = i18n.New(yaml.New("translations"))

// GetBooking returns ...
func GetBooking(c *gin.Context) {
	acceptLang := c.GetHeader("Accept-Language")

	dbc := c.MustGet("db").(*db.Client)
	fID := c.Param("id")
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	actor := c.MustGet("workflowActor").(string)
	data, state, err := dbc.GetBooking(fID, allowedDCs, actor)
	if err != nil {
		c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		return
	}
	data.State = fmt.Sprint(myI18n.T(acceptLang, data.State))
	c.JSON(http.StatusOK, *data)
}

// GetBookingPAPI returns ...
func GetBookingPAPI(c *gin.Context) {
	log.Info("GetBookingPAPI")
	dbc := c.MustGet("db").(*db.Client)
	fID := c.Param("id")
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	actor := c.MustGet("workflowActor").(string)
	data, err := dbc.GetBookingPAPI(fID, allowedDCs, actor)
	if err != nil {
		log.Errorf("Get RF fail %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if data == nil {
		log.Errorf("Not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	aar, _ := uuid.NewV4()
	url := "https:/nexcenter/v2/facility_request/{id}/file.name"
	c.JSON(
		http.StatusOK, db.BookingResponsePAPI{
			AssociatedAccessRequest: &aar,
			Customer:                &aar,
			CustomerID:              data.CustomerID,
			DataCenterID:            data.DataCenterID,
			Description:             data.Description,
			EndTime:                 data.EndTime,
			Facility:                data.Facility,
			FacilityName:            data.FacilityName,
			FileURL:                 &url,
			ID:                      data.ID,
			RequestedAt:             data.RequestedAt,
			RequestorID:             data.RequestorID,
			StartTime:               data.StartTime,
			State:                   data.State,
			StateInfo:               data.StateInfo,
			Transitions:             data.Transitions,
		},
	)
}

// GetBookings returns ...
func GetBookings(c *gin.Context) {
	log.Info("GetFacility Requests")

	dbc := c.MustGet("db").(*db.Client)
	perPage := c.MustGet("per_page").(int)
	pageNumber := c.MustGet("page_number").(int)

	request, err := makeGetRequest(c, "CAPI")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	actor := c.MustGet("workflowActor").(string)
	data, total, err := dbc.GetBookings(request, allowedDCs, actor)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input value for enum") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		log.Errorf("Error making the request: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	middlewares.WritePaginationHeaders(c, total)
	for i := range data {
		data[i].State = fmt.Sprint(myI18n.T(defaultLang, data[i].State))
	}
	c.JSON(http.StatusOK, db.BookingsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		NumResults: total,
		Objects:    data,
	})
}

// GetBookingsPAPI returns ...
func GetBookingsPAPI(c *gin.Context) {
	log.Info("GetFacility Requests")
	dbc := c.MustGet("db").(*db.Client)

	var perPage, pageNumber int
	if pp, exists := c.Get("per_page"); exists {
		perPage = pp.(int)
	}
	if p, exists := c.Get("page_number"); exists {
		pageNumber = p.(int)
	}

	request, err := makeGetRequest(c, "PAPI")
	if err != nil {
		log.Debugln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	actor := c.MustGet("workflowActor").(string)
	data, total, err := dbc.GetBookings(request, allowedDCs, actor)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	middlewares.WritePaginationHeaders(c, total)

	c.JSON(http.StatusOK, db.BookingsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		NumResults: total,
		Objects:    data,
	})
}

// PostBooking ...
func PostBooking(c *gin.Context) {
	dbc := c.MustGet("db").(*db.Client)
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)

	var body db.BookingPost

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	actor := c.MustGet("workflowActor").(string)
	response, state, err := dbc.PostBooking(&body, actor, allowedDCs)
	if err != nil {
		log.Errorf("\nError Posting FR %v \n", err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Conflict - The data you posted conflicts with the existing."})
			return
		}
		c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// PostBookingPAPI ...
func PostBookingPAPI(c *gin.Context) {
	dbc := c.MustGet("db").(*db.Client)
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	var body db.BookingPost

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Errorf("Post Facility Request failed %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	actor := c.MustGet("workflowActor").(string)
	response, state, err := dbc.PostBooking(&body, actor, allowedDCs)
	if err != nil {
		log.Errorf("Error PAPI Post FR %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Conflict - The data you posted conflicts with the existing."})
			return
		}
		c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

// PatchBooking ...
func PatchBooking(c *gin.Context) {
	dbc := c.MustGet("db").(*db.Client)
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	id := c.Param("id")
	var body db.BookingPatch

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response *db.Booking

	actor := c.MustGet("workflowActor").(string)
	response, state, err := dbc.PatchBooking(&body, id, actor, c.MustGet("UserID").(uuid.UUID).String(), allowedDCs)

	if err != nil {
		log.Errorln(err)
		if _, ok := err.(*workflow.TransitionError); ok {
			c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		}
		if strings.Contains(err.Error(), "duplicate key") {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Conflict - The data you posted conflicts with the existing."})
			return
		}
		c.AbortWithStatus(state)
		return
	}

	c.JSON(http.StatusOK, response)
}

// PatchBookingPAPI ...
func PatchBookingPAPI(c *gin.Context) {
	dbc := c.MustGet("db").(*db.Client)
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	id := c.Param("id")
	var body db.BookingPatch

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var response *db.Booking
	actor := c.MustGet("workflowActor").(string)
	response, state, err := dbc.PatchBooking(&body, id, actor, c.MustGet("UserID").(uuid.UUID).String(), allowedDCs)

	if err != nil {
		log.Errorf("\nError Patching FR %v \n", err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Conflict - The data you posted conflicts with the existing."})
			return
		}
		c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PostBookingSAPI ...
func PostBookingSAPI(c *gin.Context) {
	dbc := c.MustGet("db").(*db.Client)
	allowedDCs := c.MustGet("AuthorisedDCs").([]uuid.UUID)
	var body db.BookingPost

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	actor := c.MustGet("workflowActor").(string)
	response, state, err := dbc.PostBooking(&body, actor, allowedDCs)
	if err != nil {
		log.Errorln(err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Conflict - The data you posted conflicts with the existing."})
			return
		}
		c.AbortWithStatusJSON(state, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

// DeleteBookingSystem returns ...
func DeleteBookingSystem(c *gin.Context) {
	var err error
	log.Debugf("DeleteBooking")
	dbc := c.MustGet("db").(*db.Client)
	// no checking of access to facilities/datacenters
	fID := c.Param("id")
	//dcID := c.Query("data_center")
	dcID := c.Param("data_center")
	if dcID == "" {
		err = dbc.DeleteBookingSystem(nil, fID)
	} else {
		err = dbc.DeleteBookingSystem(&dcID, fID)
	}

	if err != nil {
		if strings.Contains(err.Error(), "pq:") {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, "Deleted OK")
}

/*======================================================================================*/
func split2ArrayString(text string) []string {
	if len(text) < 1 {
		return nil
	}
	strArr := strings.Split(text, ",")
	var response []string
	for _, s := range strArr {
		s = strings.TrimSpace(s)
		response = append(response, s)
	}
	return response
}

func split2ArrayUUID(text string) ([]uuid.UUID, error) {
	if len(text) < 1 {
		return nil, nil
	}
	strArr := strings.Split(text, ",")
	var response []uuid.UUID
	for _, s := range strArr {
		s = strings.TrimSpace(s)
		sUUID, err := uuid.FromString(s)
		if err != nil {
			return nil, err
		}
		response = append(response, sUUID)
	}
	return response, nil
}

func makeGetRequest(c *gin.Context, apiType string) (*db.BookingsGet, error) {
	var perPage, pageNumber int
	if pp, exists := c.Get("per_page"); exists {
		perPage = pp.(int)
	}
	if p, exists := c.Get("page_number"); exists {
		pageNumber = p.(int)
	}

	var customers, datacenters, facilities, requestors []uuid.UUID
	var provider uuid.UUID
	customers = c.MustGet("csList").([]uuid.UUID)
	datacenters = c.MustGet("dcList").([]uuid.UUID)
	facilities = c.MustGet("fyList").([]uuid.UUID)
	requestors = c.MustGet("rqList").([]uuid.UUID)
	if apiType == "CAPI" {
		custID, custExists := c.MustGet("customerID").(uuid.UUID)
		if custExists == false {
			return nil, errors.New("No customer-ID received")
		}
		customers = append(customers, custID)
	}

	if apiType == "PAPI" {
		custID, custExists := c.MustGet("customerID").(uuid.UUID)
		if custExists == false {
			return nil, errors.New("No customer-ID received")
		}
		provider = custID
	}

	var fromDate, toDate *time.Time
	if dt := c.Request.URL.Query().Get("fromdate"); dt != "" {
		log.Errorln("from date:", dt)
		fromD, err := time.Parse(time.RFC3339, dt)
		if err != nil {
			return nil, err
		}
		fromDate = &fromD
	}
	if dt := c.Request.URL.Query().Get("todate"); dt != "" {
		log.Errorln("to date:", dt)
		toD, err := time.Parse(time.RFC3339, dt)
		if err != nil {
			return nil, err
		}
		toDate = &toD
	}

	var states []string
	stsStr := c.Request.URL.Query().Get("states")
	states = split2ArrayString(stsStr)

	return &db.BookingsGet{
		Page:       &pageNumber,
		PerPage:    &perPage,
		DataCenter: datacenters,
		Facility:   facilities,
		Customer:   customers,
		Provider:   &provider,
		Requestor:  requestors,
		FromDate:   fromDate,
		ToDate:     toDate,
		States:     states,
	}, nil
}
