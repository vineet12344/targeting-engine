package campaign

import (
	"testing"
)

// Test MatchCampaigns against mock cache
func TestMatchCampaigns(t *testing.T) {
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
		name        string
		req         CampaignRequest
		wantCampIDs []string
	}{
		{
			name: "Match first campaign (android-IN)",
			req: CampaignRequest{
				App:     "com.test.app",
				OS:      "android",
				Country: "IN",
			},
			wantCampIDs: []string{"cmp001"},
		},
		{
			name: "Match second campaign (ios-JP)",
			req: CampaignRequest{
				App:     "com.shop.app",
				OS:      "ios",
				Country: "JP",
			},
			wantCampIDs: []string{"cmp002"},
		},
		{
			name: "No match due to excluded country",
			req: CampaignRequest{
				App:     "com.shop.app",
				OS:      "ios",
				Country: "IN",
			},
			wantCampIDs: []string{}, // Expecting no matches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchCampaigns(tt.req)

			if len(tt.wantCampIDs) == 0 {
				if len(got) > 0 {
					t.Errorf("Expected no campaigns, got %d", len(got))
				}
			} else {
				if len(got) == 0 {
					t.Errorf("Expected campaigns %v, got none", tt.wantCampIDs)
				}
				found := map[string]bool{}
				for _, c := range got {
					found[c.ID] = true
				}
				for _, expectedID := range tt.wantCampIDs {
					if !found[expectedID] {
						t.Errorf("Expected campaign %s not found in results", expectedID)
					}
				}
			}
		})
	}
}
