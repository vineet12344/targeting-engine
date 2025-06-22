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
			name: "Match all include",
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
			name: "Exclude app match fails",
			rule: TargetingRule{
				IncludeApp:     "com.test.app",
				ExcludeApp:     "com.test.app", // Exclusion overrides
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
			name: "Multiple includes match",
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
			name: "Empty include means block",
			rule: TargetingRule{
				IncludeApp: "",
			},
			req: CampaignRequest{
				App: "com.test.app",
			},
			expected: true,
		},
		{
			name: "Exclude country blocks",
			rule: TargetingRule{
				IncludeApp:     "com.test.app",
				IncludeCountry: "IN",
				ExcludeCountry: "IN", // Exclusion should take priority
			},
			req: CampaignRequest{
				App:     "com.test.app",
				Country: "IN",
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
