package campaign

import "github.com/vineet12344/targeting-engine/pkg/db"

// Fetch only Active campaigns and preload their values.
func FetchActiveCampaings() ([]Campaign, error) {
	var campaings []Campaign

	err := db.DB.Preload("Rules").Where("status=?", "ACTIVE").Find(&campaings).Error

	return campaings, err
}
