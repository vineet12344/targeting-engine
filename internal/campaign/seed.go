package campaign

import "github.com/vineet12344/targeting-engine/pkg/db"

// Seed temporary data in DB for testing
func SeedCampaings() error {
	campaigns := []Campaign{
		{
			ID:       "cmp_spotify",
			Status:   "ACTIVE",
			ImageURL: "https://ads.spotifycdn.com/download-now.png",
			CTA:      "Download",
			Rules: []TargetingRule{
				{
					IncludeCountry: "US,CA",
					// No App or OS restriction → should match all apps and OS in US or CA
				},
			},
		},
		{
			ID:       "cmp_subwaysurfer",
			Status:   "ACTIVE",
			ImageURL: "https://ads.subwaysurfer.com/play.png",
			CTA:      "Play",
			Rules: []TargetingRule{
				{
					IncludeApp:     "com.gametion.ludokinggame",
					IncludeOS:      "android",
					IncludeCountry: "US,IN",
				},
			},
		},
		{
			ID:       "cmp_discord",
			Status:   "ACTIVE",
			ImageURL: "https://ads.discord.com/join.png",
			CTA:      "Join",
			Rules: []TargetingRule{
				{
					IncludeCountry: "IN",
					IncludeOS:      "android",
				},
			},
		},
		{
			ID:       "cmp_inactive",
			Status:   "INACTIVE",
			ImageURL: "https://ads.dummy.com/inactive.png",
			CTA:      "Ignore Me",
			Rules: []TargetingRule{
				{
					IncludeCountry: "US",
					IncludeApp:     "com.gametion.ludokinggame",
					IncludeOS:      "android",
				},
			},
		},
	}

	for _, c := range campaigns {
		db.DB.Create(&c)
	}
	return nil
}
