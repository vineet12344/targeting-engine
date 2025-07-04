package campaign

// Define Database Models

type Campaign struct {
	ID       string `gorm:"primaryKey"`
	Status   string // ACTIVE or INACTIVE
	ImageURL string
	CTA      string
	Rules    []TargetingRule `gorm:"foreignKey:CampaignID"`
}

type TargetingRule struct {
	ID             uint   `gorm:"primaryKey"`
	CampaignID     string `gorm:"index"`
	IncludeApp     string
	ExcludeApp     string
	IncludeOS      string
	ExcludeOS      string
	IncludeCountry string
	ExcludeCountry string
	IncludeDevice  string
	ExcludeDevice  string
}
