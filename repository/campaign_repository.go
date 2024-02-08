package repository

import (
	"github.com/petruskuswandi/bwastartup.git/models"
	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindAllCampaign() ([]models.Campaign, error)
	FindByUserID(userID string) ([]models.Campaign, error)
	FindByID(ID string) (models.Campaign, error)
	Save(campaign models.Campaign) (models.Campaign, error)
	Update(campaign models.Campaign) (models.Campaign, error)
	CreateImage(models.CampaignImage) (models.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepositoryCampaign(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAllCampaign() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByUserID(userID string) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByID(ID string) (models.Campaign, error) {
	var campaign models.Campaign

	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Save(campaign models.Campaign) (models.Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Update(campaign models.Campaign) (models.Campaign, error) {
	err := r.db.Save(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) CreateImage(campaignImage models.CampaignImage) (models.CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(campaignID string) (bool, error) {
	err := r.db.Model(&models.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
