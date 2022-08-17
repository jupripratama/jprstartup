package campaign

import "time"

type Campaign struct {
	ID               int    `json:"id" binding:"required"`
	UserID           int    `json:"user_id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	BackerCrount     int    `json:"backer_crount" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	CurrentAmount    int    `json:"current_amount" binding:"required"`
	Slug             string `json:"slug" binding:"required"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage `json:"campaign_images" binding:"required"`
}

type CampaignImage struct {
	ID         int    `json:"id" binding:"required"`
	CampaignID int    `json:"campaign_id" binding:"required"`
	FileName   string `json:"file_name" binding:"required"`
	IsPrimary  int    `json:"isPrimary" binding:"required"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
