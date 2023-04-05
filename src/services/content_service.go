package services

import (
	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
)

type ContentService struct {
	contentRepo repository.IContentRepository
}

type IContentService interface {
	GetContent(int, string, int) (*entity.Content, error)
}

func NewContentService(contentRepo repository.IContentRepository) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
	}
}

func (s *ContentService) GetContent(serviceId int, name string, pin int) (*entity.Content, error) {
	result, err := s.contentRepo.Get(serviceId, name)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if result != nil {
		content = entity.Content{
			Value: result.Value,
			Tid:   result.Tid,
		}

		// set pin
		if pin > 0 {
			content.SetPIN(pin)
		}
	}
	return &content, nil
}
