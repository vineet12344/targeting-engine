package campaign

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vineet12344/targeting-engine/middleware"
)

type CampaignResponse struct {
	ID string `json:"ID"`
	// Status   string `json:"Status"`
	ImageURL string `json:"ImageURL"`
	CTA      string `json:"CTA"`
}

type CampaignRequest struct {
	App     string
	OS      string
	Country string
	Device  string
}

// Return matched campaigns accordign to user  Query
func RegisterRoutes(router *gin.Engine) {

	router.GET("/campaign", func(c *gin.Context) {
		app := c.Query("app")
		os := c.Query("os")
		country := c.Query("country")
		device := c.Query("device")

		if app == "" || os == "" || country == "" || device == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query params"})
			return
		}

		req := CampaignRequest{
			App:     app,
			OS:      os,
			Country: country,
			Device:  device,
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

		var responses []CampaignResponse
		for _, c := range matched {
			responses = append(responses, CampaignResponse{
				ID:       c.ID,
				ImageURL: c.ImageURL,
				CTA:      c.CTA,
				// Status:   c.Status,
			})
		}

		c.JSON(http.StatusOK, responses)

	})
}
