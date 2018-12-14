package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Pagination takes care of the pagination parameters and headers
func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		var perPage, pageNumber int
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodPost {

			str, found := c.GetQuery("per_page")
			if found {
				n, err := strconv.ParseInt(str, 10, 32)
				if err != nil || n <= 0 {
					e := fmt.Sprintf("Invalid value %q - expected an integer greater than zero", str)
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Validation error",
						"details": gin.H{
							"per_page": e,
						},
					})
					c.Abort()
					return
				}
				perPage = int(n)
			} else {
				perPage = 100
			}

			str, found = c.GetQuery("page")
			if found {
				n, err := strconv.ParseInt(str, 10, 32)
				if err != nil || n <= 0 {
					e := fmt.Sprintf("Invalid value %q - expected an integer greater than zero", str)
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Validation error",
						"details": gin.H{
							"page": e,
						},
					})
					c.Abort()
					return
				}
				pageNumber = int(n)
			} else {
				pageNumber = 1
			}

			c.Set("page_number", pageNumber)
			c.Set("per_page", perPage)
		}
		c.Next()
	}
}

// WritePaginationHeaders sets the response pagination headers
func WritePaginationHeaders(c *gin.Context, totalCount int) {
	pageNumber := c.MustGet("page_number").(int)
	perPage := c.MustGet("per_page").(int)
	c.Header("X-Page", fmt.Sprintf("%d", pageNumber))
	c.Header("X-Per-Page", fmt.Sprintf("%d", perPage))
	c.Header("X-Total-Count", fmt.Sprintf("%d", totalCount))
	c.Header("Link", createLinkHeader(pageNumber, perPage, totalCount))
}

func createLinkHeader(pageNumber, perPage, totalCount int) string {
	ret := []string{}
	if pageNumber != 1 {
		ret = append(ret, fmt.Sprintf(`?page=%d&per_page=%d, rel="first"`, 1, perPage))
		ret = append(ret, fmt.Sprintf(`?page=%d&per_page=%d, rel="prev"`, pageNumber-1, perPage))
	}
	if pageNumber*perPage < totalCount { // not the last page
		ret = append(ret, fmt.Sprintf(`?page=%d&per_page=%d, rel="next"`, pageNumber+1, perPage))
		lastPage := totalCount / perPage
		if lastPage*perPage < totalCount {
			lastPage++ // we need a ceiling division, integer division is by default floored; if lastPage*perPage == totalCount they're the same, otherwise we need to add one
		}
		ret = append(ret, fmt.Sprintf(`?page=%d&per_page=%d, rel="last"`, lastPage, perPage))
	}
	return strings.Join(ret, "; ")
}
