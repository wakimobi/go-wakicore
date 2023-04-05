package services

import (
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
)

type TransactionService struct {
	transactionRepo repository.ITransactionRepository
}

type ITransactionService interface {
	SaveTransaction(*entity.Transaction) error
	UpdateTransaction(*entity.Transaction) error
}

func NewTransactionService(transactionRepo repository.ITransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) SaveTransaction(t *entity.Transaction) error {
	err := s.transactionRepo.Save(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) UpdateTransaction(t *entity.Transaction) error {
	data := &entity.Transaction{
		ServiceID: t.ServiceID,
		Msisdn:    t.Msisdn,
		Subject:   t.Subject,
		Status:    "FAILED",
	}
	errDelete := s.transactionRepo.Delete(data)
	if errDelete != nil {
		return errDelete
	}

	errSave := s.transactionRepo.Save(t)
	if errSave != nil {
		return errSave
	}
	return nil
}
