package services

import (
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
)

type CampaignService struct {
	campaignRepo repository.ICampaignRepository
}

type ICampaignService interface {
	GetCampaign(int, string) bool
	UpdateCampaign(int, string) error
}

func NewCampaignService(campaignRepo repository.ICampaignRepository) *CampaignService {
	return &CampaignService{
		campaignRepo: campaignRepo,
	}
}

func (s *CampaignService) GetCampaign(serviceId int, adnet string) bool {
	count, _ := s.campaignRepo.Count(serviceId, adnet)
	return count > 0
}

func (s *CampaignService) UpdateCampaign(serviceId int, adnet string) error {
	err := s.campaignRepo.Update(
		&entity.Campaign{
			ServiceID: serviceId,
			Adnet:     adnet,
		})

	if err != nil {
		return err
	}
	return nil
}
