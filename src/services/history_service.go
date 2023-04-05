package services

import (
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
)

type HistoryService struct {
	transactionRepo repository.IHistoryRepository
}

type IHistoryService interface {
	SaveHistory(*entity.History) error
}

func NewHistoryService(transactionRepo repository.IHistoryRepository) *HistoryService {
	return &HistoryService{
		transactionRepo: transactionRepo,
	}
}

func (s *HistoryService) SaveHistory(t *entity.History) error {
	err := s.transactionRepo.Save(t)
	if err != nil {
		return err
	}
	return nil
}
