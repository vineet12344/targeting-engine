package campaign

import (
	"log"
	"strings"
	"sync"
)

type CampaignRequest struct {
	App     string
	OS      string
	Country string
}

// Main Logic
func MatchCampaigns(req CampaignRequest) []Campaign {
	campaigns := GetCachedCampaigns()
	// var matches []Campaign
	// log.Println("🔍 Matching for request:", req)

	// for _, c := range campaigns {
	// 	for _, rule := range c.Rules {
	// 		if ruleMatches(rule, req) {
	// 			matches = append(matches, c)
	// 			break
	// 		}
	// 	}
	// }

	resultChan := make(chan Campaign, len(campaigns))
	var wg sync.WaitGroup
	var matches []Campaign

	for i, c := range campaigns {
		campaign := c
		wg.Add(1)

		go func(i int) {
			log.Print("Running go-routine: ", i)
			defer wg.Done()

			for _, rule := range campaign.Rules {
				if ruleMatches(rule, req) {
					resultChan <- campaign
					break
				}
			}
		}(i)

	}

	wg.Wait()
	close(resultChan)

	for c := range resultChan {
		matches = append(matches, c)
	}

	return matches
}

func ruleMatches(rule TargetingRule, req CampaignRequest) bool {
	// Match Include values(if present)
	if !matchesInclude(rule.IncludeApp, req.App) {
		return false
	}
	if !matchesInclude(rule.IncludeCountry, req.Country) {
		return false
	}
	if !matchesInclude(rule.IncludeOS, req.OS) {
		return false
	}

	// Match Exclude values
	if matchesExclude(rule.ExcludeApp, req.App) {
		return false
	}

	if matchesExclude(rule.ExcludeCountry, req.Country) {
		return false
	}

	if matchesExclude(rule.ExcludeOS, req.OS) {
		return false
	}

	return true

}

func matchesInclude(list string, target string) bool {
	if list == "" {
		return true //No restrictions
	}

	values := strings.Split(list, ",")
	for _, v := range values {
		if strings.TrimSpace(v) == target {
			return true
		}
	}

	return false

}

func matchesExclude(list string, target string) bool {
	if list == "" {
		return false // nothing excluded
	}
	values := strings.Split(list, ",")
	for _, v := range values {
		if strings.TrimSpace(v) == target {
			return true
		}
	}
	return false
}
