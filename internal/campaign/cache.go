package campaign

import (
	"log"
	"sync"
	"time"
)

var (
	cachedCampaigns []Campaign
	cacheMutex      sync.RWMutex
)

func LoadToCache(service CampaignService) error {
	campaigns, err := service.GetActiveCampaings()
	if err != nil {
		return err
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cachedCampaigns = campaigns
	log.Println("✅ Cache refreshed at", time.Now())
	// log.Println("📦Current Cache is :", cachedCampaigns)

	return nil
}

func GetCachedCampaigns() []Campaign {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return cachedCampaigns
}

func StartAutoRefresh(service CampaignService, interval time.Duration, stopChan <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				LoadToCache(service)
			case <-stopChan:
				return
			}
		}
	}()
}
