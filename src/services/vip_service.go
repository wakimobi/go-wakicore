package services

import "github.com/idprm/go-pass-tsel/src/domain/repository"

type VIPService struct {
	vipRepo repository.IVIPRepository
}

type IVIPService interface {
	GetVIP(string) bool
}

func NewVIPService(vipRepo repository.IVIPRepository) *VIPService {
	return &VIPService{
		vipRepo: vipRepo,
	}
}

func (s *VIPService) GetVIP(msisdn string) bool {
	count, _ := s.vipRepo.Count(msisdn)
	return count > 0
}
