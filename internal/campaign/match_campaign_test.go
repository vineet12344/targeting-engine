package campaign

import (
	"testing"
)

// Test MatchCampaign against mock cache
func TestMatchCampaign(t *testing.T) {
	mockCampaigns := []Campaign{
		{
			ID:       "cmp001",
			Status:   "ACTIVE",
			ImageURL: "https://test.com/ad1.png",
			CTA:      "Download Now",
			Rules: []TargetingRule{
				{
					IncludeApp:     "com.test.app",
					IncludeOS:      "android",
					IncludeCountry: "IN,US",
				},
			},
		},
		{
			ID:       "cmp002",
			Status:   "ACTIVE",
			ImageURL: "https://test.com/ad2.png",
			CTA:      "Try Now",
			Rules: []TargetingRule{
				{
					IncludeApp:     "com.shop.app",
					IncludeOS:      "ios",
					ExcludeCountry: "IN",
				},
			},
		},
	}

	// Set mock campaigns to in-memory cache
	SetCachedCampaigns(mockCampaigns)

	tests := []struct {
		name       string
		req        CampaignRequest
		wantCampID string
	}{
		{
			name: "Match first campaign (android-IN)",
			req: CampaignRequest{
				App:     "com.test.app",
				OS:      "android",
				Country: "IN",
			},
			wantCampID: "cmp001",
		},
		{
			name: "Match second campaign (ios-JP)",
			req: CampaignRequest{
				App:     "com.shop.app",
				OS:      "ios",
				Country: "JP",
			},
			wantCampID: "cmp002",
		},
		{
			name: "No match due to excluded country",
			req: CampaignRequest{
				App:     "com.shop.app",
				OS:      "ios",
				Country: "IN",
			},
			wantCampID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchCampaign(tt.req)

			if tt.wantCampID == "" && got != nil {
				t.Errorf("Expected nil campaign, got %+v", got)
			}

			if tt.wantCampID != "" {
				if got == nil {
					t.Errorf("Expected campaign %s, got nil", tt.wantCampID)
				} else if got.ID != tt.wantCampID {
					t.Errorf("Expected campaign %s, got %s", tt.wantCampID, got.ID)
				}
			}
		})
	}
}
