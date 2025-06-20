package campaign

type Campaign struct {
	ID       string `gorm:"primarykey"`
	Status   string
	ImageURL string
	CTA      string
	Rules    []TargetingRule `gorm:"foreignKey:CampaignID"`
}

type TargetingRule struct {
	ID             uint `gorm:"primaryKey"`
	CampaignID     string
	IncludeApp     string
	ExcludeApp     string
	IncludeOS      string
	ExcludeOS      string
	IncludeCountry string
	ExcludeCountry string
}
