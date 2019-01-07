package handlers

import (
	"bookings/dbmodels"
	"fmt"
	"net/http"
	"strings"

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
	bID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	actor := c.MustGet("workflowActor").(string)
	data, httpstatus, err := dbmodels.GetBooking(bID, actor)
	if err != nil {
		c.AbortWithStatusJSON(httpstatus, gin.H{"message": err.Error()})
		return
	}
	data.State = fmt.Sprint(myI18n.T(acceptLang, data.State))
	c.JSON(http.StatusOK, *data)
}

// GetBookings returns ...
func GetBookings(c *gin.Context) {
	log.Info("GetBookings Requests")

	perPage := c.MustGet("per_page").(int)
	pageNumber := c.MustGet("page_number").(int)

	actor := c.MustGet("workflowActor").(string)
	data, total, err := dbmodels.GetBookings(actor)
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
	c.JSON(http.StatusOK, dbmodels.BookingsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		NumResults: total,
		Objects:    data,
	})
}

// GetBookingsPAPI returns ...
func GetBookingsPAPI(c *gin.Context) {
	log.Info("GetBookingsPAPI Requests")

	var perPage, pageNumber int
	if pp, exists := c.Get("per_page"); exists {
		perPage = pp.(int)
	}
	if p, exists := c.Get("page_number"); exists {
		pageNumber = p.(int)
	}

	actor := c.MustGet("workflowActor").(string)
	data, total, err := dbmodels.GetBookings(actor)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	middlewares.WritePaginationHeaders(c, total)

	c.JSON(http.StatusOK, dbmodels.BookingsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		NumResults: total,
		Objects:    data,
	})
}

// PostBooking ...
func PostBooking(c *gin.Context) {
	var body dbmodels.BookingPost
	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	response, state, err := dbmodels.PostBooking(&body)
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
	var body dbmodels.BookingPost
	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Errorf("Post PostBookingPAPI Request failed %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	response, state, err := dbmodels.PostBooking(&body)
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
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad booking ID"})
		return
	}
	var body dbmodels.BookingPatch

	err = c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response *dbmodels.Booking

	response, state, err := dbmodels.PatchBooking(&body, id)

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
	id, err := uuid.FromString(c.Param("id"))
	var body dbmodels.BookingPatch

	err = c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var response *dbmodels.Booking

	response, state, err := dbmodels.PatchBooking(&body, id)

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
	var body dbmodels.BookingPost

	err := c.MustBindWith(&body, binding.JSON)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	body.CustomerID = c.MustGet("customerID").(uuid.UUID)
	body.RequestorID = c.MustGet("UserID").(uuid.UUID)

	response, state, err := dbmodels.PostBooking(&body)
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
	bID, err := uuid.FromString(c.Param("id"))
	hotelID, err := uuid.FromString(c.Param("hotelid"))
	if hotelID == uuid.Nil {
		err = dbmodels.DeleteBookingSystem(nil, bID)
	} else {
		err = dbmodels.DeleteBookingSystem(&hotelID, bID)
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
