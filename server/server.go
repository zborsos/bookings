package server

import (
	"bookings/middleware"
	"net/http"
	"regexp"
	"strings"

	"github.com/miketonks/swag"
	sv "github.com/miketonks/swag-validator"
	"github.com/miketonks/swag/swagger"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

const commVersion = "0.0.0.1"

// checkHeaders checks x-org-header is set
func checkHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		endPointURL := c.Request.URL.String()
		userIDstr := c.Request.Header.Get("X-Customer")
		//IF YOU NEED TO SEE THE HEADER
		for name, headers := range c.Request.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				log.Debugf("HEADER data - %v:\t%v\n", name, h)
			}
		}

		if userIDstr == "" {
			c.AbortWithStatus(403)
			return
		}
		userID, err := uuid.FromString(userIDstr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid user ID"})
		}
		if userID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid user ID NIL user is not accepted"})
		}

		c.Set("userID", userID)

		isDocReq, _ := regexp.MatchString("/bookings/booking-doc", endPointURL)
		if !isDocReq {
			match, err := regexp.MatchString("/bookings/(provider|system)", endPointURL)
			actor := "Customer"
			if err != nil {
				c.AbortWithStatus(500)
			}

			if match {
				actor = "Provider"
			}
			c.Set("workflowActor", actor)
		}

		c.Next()
	}
}

// RunServer runs the server
func RunServer() {
	r := CreateRouter()
	err := r.Run(":5670")
	if err != nil {
		log.Fatalf("server exited: %s", err)
	}
}

// CreateRouter creates the router
func CreateRouter() *gin.Engine {

	r := gin.New()

	// set context objects
	r.Use(func(c *gin.Context) {

		c.Next()
	})

	capi := CreateSwaggerCAPI()
	papi := CreateSwaggerPAPI()
	sapi := CreateSwaggerSAPI()
	enableCors := false

	org := r.Group("", checkHeaders(), sv.SwaggerValidator(capi), sv.SwaggerValidator(papi), sv.SwaggerValidator(sapi), middleware.Pagination(), middleware.ValidateUUIDs())

	org.GET("/bookings/bookings-doc", gin.WrapH(capi.Handler(enableCors)))
	org.GET("/bookings/provider/bookings-doc", gin.WrapH(papi.Handler(enableCors)))
	org.GET("/bookings/system/bookings-doc", gin.WrapH(sapi.Handler(enableCors)))

	capi.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)
		/*
			if endpoint.Method == "POST" || endpoint.Method == "PATCH" {
				org.Handle(endpoint.Method, path, middleware, h)
			}
		*/
		org.Handle(endpoint.Method, path, h)
	})
	papi.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)
		/*
			if endpoint.Method == "POST" || endpoint.Method == "PATCH" {
				org.Handle(endpoint.Method, path, middleware, h)
			}
		*/
		org.Handle(endpoint.Method, path, h)
	})
	sapi.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)
		/*
			if endpoint.Method == "POST" || endpoint.Method == "PATCH" {
				org.Handle(endpoint.Method, path, middleware, h)
			}
		*/
		org.Handle(endpoint.Method, path, h)
	})
	return r
}

// CreateSwaggerCAPI ...
func CreateSwaggerCAPI() *swagger.API {
	api := swag.New(
		swag.Title("BOOKINGS CUSTOMER API"),
		swag.Description("This API manages user bookings"),
		swag.Version(commVersion),
		swag.BasePath("/bookings"),
		swag.Endpoints(
			aggregateEndpoints(
				bookingsCAPI(),
			)...,
		),
	)

	return api
}

// CreateSwaggerPAPI ...
func CreateSwaggerPAPI() *swagger.API {
	api := swag.New(
		swag.Title("BOOKINGS PARTNER API"),
		swag.Description("This API manages partners bookings"),
		swag.Version(commVersion),
		swag.BasePath("/bookings"),
		swag.Endpoints(
			aggregateEndpoints(
				bookingsPAPI(),
			)...,
		),
	)

	return api
}

// CreateSwaggerSAPI ...
func CreateSwaggerSAPI() *swagger.API {
	api := swag.New(
		swag.Title("BOOKINGS SYSYTEM API"),
		swag.Description("This API is for Sysytem use only"),
		swag.Version(commVersion),
		swag.BasePath("/bookings"),
		swag.Endpoints(
			aggregateEndpoints(
				bookingsSAPI(),
			)...,
		),
	)

	return api
}

func aggregateEndpoints(endpoints ...[]*swagger.Endpoint) []*swagger.Endpoint {
	res := []*swagger.Endpoint{}
	for _, v := range endpoints {
		res = append(res, v...)
	}
	return res
}
