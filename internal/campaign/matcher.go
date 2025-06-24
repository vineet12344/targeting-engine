package campaign

import (
	"log"
	"strings"
	"sync"
)

func MatchBatchCampaigns(req CampaignRequest) []Campaign {
	campaigns := GetCachedCampaigns()
	numWorkers := 10

	type Result struct {
		match Campaign
	}

	jobs := make(chan Campaign, len(campaigns))
	result := make(chan Result, len(campaigns))

	var wg sync.WaitGroup

	// Here we start N workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(w int) {
			log.Printf("Running Go-routine %v", w)
			defer wg.Done()
			for campaign := range jobs {
				for _, rule := range campaign.Rules {
					if ruleMatches(rule, req) {
						result <- Result{match: campaign}
						break
					}
				}
			}
		}(w)
	}

	for _, c := range campaigns {
		jobs <- c
	}

	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	close(result)

	var matches []Campaign
	seen := make(map[string]bool)

	for res := range result {
		if !seen[res.match.ID] {
			matches = append(matches, res.match)
			seen[res.match.ID] = true
		}
	}

	return matches

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

	// wait for all doroutines to complete their jobs
	wg.Wait()
	// close result channel
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
	if !matchesInclude(rule.IncludeOS, req.OS) {
		return false
	}

	if !matchesInclude(rule.IncludeDevice, req.Device) {
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

	if matchesExclude(rule.ExcludeDevice, req.Device) {
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
