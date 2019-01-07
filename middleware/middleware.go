package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// ValidateUUIDs validates the list of DCs in the query
func ValidateUUIDs() gin.HandlerFunc {
	fields := map[string]struct {
		ContextName string
		MessageName string
	}{
		"data_center": {
			ContextName: "dcList",
			MessageName: "Data Center",
		},
		"customer": {
			ContextName: "csList",
			MessageName: "Customer",
		},
		"room": {
			ContextName: "roomList",
			MessageName: "Room",
		},
		"requestor": {
			ContextName: "rqList",
			MessageName: "Requestor",
		},
	}

	return func(c *gin.Context) {
		for k, field := range fields {
			IDs, exists := c.GetQuery(k)
			var UUIDs []uuid.UUID
			if exists {
				items := strings.Split(IDs, ",")
				for _, v := range items {
					item := strings.TrimSpace(v)
					uuid, err := uuid.FromString(item)
					if err != nil {
						msg := fmt.Sprintf("%s ->%s<- is Not valid UUID", field.MessageName, item)
						log.Error(msg)
						c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
						return
					}
					UUIDs = append(UUIDs, uuid)
				}
			}
			c.Set(field.ContextName, UUIDs)
		}
	}
}

/*
var (
	uuid2id     = make(map[string]int)
	param2table = map[string]string{
		"dc_id":      "data_centers",
		"pg_id":      "colo_permission_groups",
		"user":       "contacts",
		"cab":        "colo_cabs",
		"permission": "dc_colo_permissions",
	}
)

// CheckHeaders checks customer and user glass headers
func CheckHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbc := c.MustGet("dbc").(*db.Client)
		glass := c.MustGet("glass").(glass.Interface)

		id := c.Request.Header.Get("Glass-Customer-id")
		if id == "" {
			msg := "Glass-Customer-id header is missing"
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}
		glassCustID, err := strconv.Atoi(id)
		if err != nil {
			msg := fmt.Sprintf("Invalid customer glass id: %v", id)
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}

		custID := dbc.GetIDByGlassID("customers", glassCustID)
		if custID == 0 {
			msg := fmt.Sprintf("Customer glass id not found: %v", glassCustID)
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}

		id = c.Request.Header.Get("Glass-User-id")
		if id == "" {
			msg := "Glass-User-id header is missing"
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}
		glassUserID, err := strconv.Atoi(id)
		if err != nil {
			msg := fmt.Sprintf("Invalid user glass id: %s", id)
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}

		userID := dbc.GetIDByGlassID("contacts", glassUserID)
		if userID == 0 {
			msg := fmt.Sprintf("User glass id not found: %v", glassUserID)
			log.Error(msg)
			c.AbortWithStatusJSON(400, gin.H{"code": 400, "message": msg})
			return
		}

		personID := dbc.GetPersonID(userID)

		user, err := glass.GetUser(strconv.Itoa(glassCustID), strconv.Itoa(glassUserID))
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(500)
			return
		}
		if user.PortalAdmin == false {
			log.Errorf("User %s is not a portal admin", *user.Login)
			c.AbortWithStatusJSON(403, gin.H{"code": 403, "message": "Forbidden"})
			return
		}

		c.Set("glassCustID", glassCustID)
		c.Set("custID", custID)
		c.Set("glassUserID", glassUserID)
		c.Set("userID", userID)
		c.Set("personID", personID)

		c.Next()

	}
}
*/
