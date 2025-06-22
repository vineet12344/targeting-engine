package campaign

import "testing"

func TestRuleMatches(t *testing.T) {
	tests := []struct {
		name     string
		rule     TargetingRule
		req      CampaignRequest
		expected bool
	}{
		{
			name: "Match all includes",
			rule: TargetingRule{
				IncludeApp:     "com.test.app",
				IncludeOS:      "android",
				IncludeCountry: "IN",
			},
			req: CampaignRequest{
				App:     "com.test.app",
				OS:      "android",
				Country: "IN",
			},
			expected: true,
		},
		{
			name: "Excluded app overrides include",
			rule: TargetingRule{
				IncludeApp:     "com.test.app",
				ExcludeApp:     "com.test.app", // exclusion should override
				IncludeOS:      "android",
				IncludeCountry: "IN",
			},
			req: CampaignRequest{
				App:     "com.test.app",
				OS:      "android",
				Country: "IN",
			},
			expected: false,
		},
		{
			name: "Multiple includes matched",
			rule: TargetingRule{
				IncludeApp:     "com.shop.app,com.test.app",
				IncludeOS:      "android,ios",
				IncludeCountry: "IN,JP",
			},
			req: CampaignRequest{
				App:     "com.test.app",
				OS:      "ios",
				Country: "JP",
			},
			expected: true,
		},
		{
			name: "Empty include means allow",
			rule: TargetingRule{
				IncludeApp: "",
			},
			req: CampaignRequest{
				App: "com.test.app",
			},
			expected: true, // ✅ If empty include means no restriction
		},
		{
			name: "Exclude country blocks",
			rule: TargetingRule{
				IncludeApp:     "com.test.app",
				IncludeCountry: "IN",
				ExcludeCountry: "IN", // Exclude overrides include
			},
			req: CampaignRequest{
				App:     "com.test.app",
				Country: "IN",
			},
			expected: false,
		},
		{
			name: "ExcludeCountry excludes even if Include matches",
			rule: TargetingRule{
				IncludeCountry: "US,IN",
				ExcludeCountry: "US",
			},
			req: CampaignRequest{
				Country: "US",
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ruleMatches(tc.rule, tc.req)
			if result != tc.expected {
				t.Errorf("Test '%s' failed. Expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}
