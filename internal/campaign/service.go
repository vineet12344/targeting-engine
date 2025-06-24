package campaign

// Methods under CampaignbService interface
type CampaignService interface {
	GetActiveCampaings() ([]Campaign, error)
}

type campaignService struct{}

func NewCampaignService() CampaignService {
	return &campaignService{}
}

func (c *campaignService) GetActiveCampaings() ([]Campaign, error) {
	return FetchActiveCampaings()
}
