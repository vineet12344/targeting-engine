package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/campaign", func(c *gin.Context) {
		app := c.Query("app")
		os := c.Query("os")
		country := c.Query("country")

		if app == "" || os == "" || country == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query params"})
			return
		}

		req := CampaignRequest{
			App:     app,
			OS:      os,
			Country: country,
		}

		matched := MatchCampaign(req)

		if matched == nil {
			c.JSON(http.StatusOK, gin.H{"message": "No matching campaign found"})

			return
		}

		c.JSON(http.StatusOK, matched)

	})
}
