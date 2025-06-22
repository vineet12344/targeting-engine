package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vineet12344/targeting-engine/middleware"
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

		matched := MatchCampaigns(req)

		if len(matched) == 0 {
			middleware.MatchMissCount.Inc()
			middleware.LogMetricValue("Match_Miss_Count", middleware.MatchMissCount)
			c.JSON(http.StatusNoContent, gin.H{"message": "No matching campaign found"})

			return
		}

		middleware.MatchSuccessCount.Inc()
		middleware.LogMetricValue("Match_SUCCESS_Count", middleware.MatchSuccessCount)

		c.JSON(http.StatusOK, matched)

	})
}
