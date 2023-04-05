package services

import (
	"log"

	"github.com/idprm/go-pass-tsel/src/domain/entity"
	"github.com/idprm/go-pass-tsel/src/domain/repository"
)

type SubscriptionService struct {
	subscriptionRepo repository.ISubscriptionRepository
}

type ISubscriptionService interface {
	GetActiveSubscription(int, string) bool
	GetSubscription(int, string) bool
	SaveSubscription(*entity.Subscription) error
	UpdateSuccess(*entity.Subscription) error
	UpdateFailed(*entity.Subscription) error
	UpdateLatest(*entity.Subscription) error
	UpdateEnable(*entity.Subscription) error
	UpdateDisable(*entity.Subscription) error
	ReminderSubscription() *[]entity.Subscription
	RenewalSubscription() *[]entity.Subscription
	RetrySubscription() *[]entity.Subscription
	AveragePerUser() *[]entity.AveragePerUser
}

func NewSubscriptionService(subscriptionRepo repository.ISubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionService) GetActiveSubscription(serviceId int, msisdn string) bool {
	count, err := s.subscriptionRepo.CountActive(serviceId, msisdn)
	if err != nil {
		log.Println(err)
	}
	return count > 0
}

func (s *SubscriptionService) GetSubscription(serviceId int, msisdn string) bool {
	count, err := s.subscriptionRepo.Count(serviceId, msisdn)
	if err != nil {
		log.Println(err)
	}
	return count > 0
}

func (s *SubscriptionService) SaveSubscription(sub *entity.Subscription) error {
	err := s.subscriptionRepo.Save(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateSuccess(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateSuccess(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateFailed(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateFailed(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateLatest(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateLatest(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateEnable(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateEnable(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateDisable(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateDisable(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) ReminderSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Reminder()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RenewalSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Renewal()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RetrySubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Retry()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) AveragePerUser() *[]entity.AveragePerUser {
	subs, err := s.subscriptionRepo.AveragePerUser("20", "20", "20")
	if err != nil {
		log.Println(err)
	}
	return subs
}
