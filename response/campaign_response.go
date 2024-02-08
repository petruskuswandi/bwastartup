package response

import (
	"strings"

	"github.com/petruskuswandi/bwastartup.git/models"
)

type CampaignResponse struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func ResponseCampaign(campaign models.Campaign) CampaignResponse {
	campaignFormatter := CampaignResponse{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func ResponseCampaigns(campaigns []models.Campaign) []CampaignResponse {
	campaignsFormatter := []CampaignResponse{}

	for _, campaign := range campaigns {
		campaignFormatter := ResponseCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

type CampaignDetailResponse struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	ImageURL         string                  `json:"image_url"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	BackerCount      int                     `json:"backer_count"`
	UserID           string                  `json:"user_id"`
	Slug             string                  `json:"slug"`
	Perks            []string                `json:"perks"`
	User             CampaignUserResponse    `json:"user"`
	Images           []CampaignImageResponse `json:"images"`
}

type CampaignUserResponse struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageResponse struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func ResponseCampaignDetail(campaign models.Campaign) CampaignDetailResponse {
	campaignDetailFormatter := CampaignDetailResponse{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.BackerCount = campaign.BackerCount
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string

	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}

	campaignDetailFormatter.Perks = perks

	user := campaign.User

	campaignUserResponse := CampaignUserResponse{}
	campaignUserResponse.Name = user.Name
	campaignUserResponse.ImageURL = user.AvatarFileName

	campaignDetailFormatter.User = campaignUserResponse

	images := []CampaignImageResponse{}

	for _, image := range campaign.CampaignImages {
		CampaignImageResponse := CampaignImageResponse{}
		CampaignImageResponse.ImageURL = image.FileName

		isPrimary := false

		if image.IsPrimary == 1 {
			isPrimary = true
		}
		CampaignImageResponse.IsPrimary = isPrimary

		images = append(images, CampaignImageResponse)
	}

	campaignDetailFormatter.Images = images

	return campaignDetailFormatter
}
