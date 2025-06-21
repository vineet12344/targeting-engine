package campaign

import "github.com/vineet12344/targeting-engine/pkg/db"

func SeedCampaings() error {
	campaings := []Campaign{
		{
			ID:       "cmp123",
			Status:   "ACTIVE",
			ImageURL: "https://example.com/ad1.png",
			CTA:      "Install Now",
			Rules: []TargetingRule{
				{
					IncludeCountry: "IN,JP,US",
					IncludeApp:     "com.test.app",
					IncludeOS:      "android",
				},
			},
		},
		{
			ID:       "cmp456",
			Status:   "ACTIVE",
			ImageURL: "https://example.com/ad2.png",
			CTA:      "Try Free",
			Rules: []TargetingRule{
				{
					ExcludeCountry: "US,IN",
					IncludeApp:     "com.shop.app",
					IncludeOS:      "ios",
				},
			},
		},
	}

	for _, c := range campaings {
		db.DB.Create(&c)
	}

	return nil

}
