package services

import (
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
)

type IDashboardService interface {
	Dashboard(startDate string, endDate string) (*models.ResponseDashboard, error)
}
type DashboardService struct {
	repo repositories.IDashboardRepository
}

func (s *DashboardService) Dashboard(startDate string, endDate string) (*models.ResponseDashboard, error) {
	return s.repo.Dashboard(startDate, endDate)
}

func NewDashboardService(repo repositories.IDashboardRepository) IDashboardService {
	return &DashboardService{repo: repo}
}
