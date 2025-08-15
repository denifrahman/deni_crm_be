package services

import (
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
)

type ISubscriptionsService interface {
	CreateSubscription(subscription *models.SubscriptionRequestCreate) error
	GetAllSubscriptions(page, size int, startDate, endDate, search string) (*models.SubscriptionResponse, error)
	ExportSubscriptionToExcel(startDate, endDate, search string) (*models.SubscriptionResponse, error)
	GetSubscriptionByID(id uint) (*models.Subscription, error)
	UpdateSubscription(subscription *models.SubscriptionRequestUpdate) error
	DeleteSubscription(id uint) error
}

type SubscriptionsService struct {
	repo repositories.ISubscriptionsRepository
}

func NewSubscriptionsService(repo repositories.ISubscriptionsRepository) ISubscriptionsService {
	return &SubscriptionsService{repo: repo}
}

func (s *SubscriptionsService) CreateSubscription(subscription *models.SubscriptionRequestCreate) error {
	return s.repo.Create(subscription.ToModel())
}

func (s *SubscriptionsService) GetAllSubscriptions(page, size int, startDate, endDate, search string) (*models.SubscriptionResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.SubscriptionToResponse(models.SubscriptionResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *SubscriptionsService) ExportSubscriptionToExcel(startDate, endDate, search string) (*models.SubscriptionResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.SubscriptionToResponse(models.SubscriptionResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *SubscriptionsService) GetSubscriptionByID(id uint) (*models.Subscription, error) {
	return s.repo.FindByID(id)
}

func (s *SubscriptionsService) UpdateSubscription(subscription *models.SubscriptionRequestUpdate) error {
	return s.repo.Update(subscription.ToModel())
}

func (s *SubscriptionsService) DeleteSubscription(id uint) error {
	return s.repo.Delete(id)
}
