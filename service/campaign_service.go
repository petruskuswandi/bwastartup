package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/petruskuswandi/bwastartup.git/models"
	"github.com/petruskuswandi/bwastartup.git/repository"
	"github.com/petruskuswandi/bwastartup.git/request"
)

type ServiceCampaign interface {
	GetCampaigns(userID string) ([]models.Campaign, error)
	GetCampaignByID(input request.GetCampaignDetailInput) (models.Campaign, error)
	CreateCampaign(input request.CreateCampaignInput) (models.Campaign, error)
	UpdateCampaign(inputID request.GetCampaignDetailInput, inputData request.CreateCampaignInput) (models.Campaign, error)
	SaveCampaignImage(input request.CreateCampaignImageInput, fileLocation string) (models.CampaignImage, error)
}

type service struct {
	repository repository.CampaignRepository
}

func NewServiceCampaign(repo repository.CampaignRepository) *service {
	return &service{repo}
}

func (s *service) GetCampaigns(UserID string) ([]models.Campaign, error) {
	if UserID != "" {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAllCampaign()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input request.GetCampaignDetailInput) (models.Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input request.CreateCampaignInput) (models.Campaign, error) {
	campaign := models.Campaign{}
	campaignID, err := uuid.NewRandom()
	if err != nil {
		return campaign, err
	}
	campaign.ID = campaignID.String()
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %s", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UpdateCampaign(inputID request.GetCampaignDetailInput, inputData request.CreateCampaignInput) (models.Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input request.CreateCampaignImageInput, fileLocation string) (models.CampaignImage, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return models.CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return models.CampaignImage{}, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return models.CampaignImage{}, err
		}
	}

	campaignImageID, err := uuid.NewRandom()
	if err != nil {
		return models.CampaignImage{}, err
	}
	campaignImage := models.CampaignImage{}
	campaignImage.ID = campaignImageID.String()
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaign, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}
