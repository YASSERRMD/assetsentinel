package services

import (
	"assetsentinel/internal/repository"
)

type MaintenanceService struct {
	repo *repository.Repository
}

func NewMaintenanceService(repo *repository.Repository) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) Create(plan *repository.MaintenancePlan) error {
	return s.repo.CreateMaintenancePlan(plan)
}

func (s *MaintenanceService) Get(id, orgID uint) (*repository.MaintenancePlan, error) {
	return s.repo.GetMaintenancePlan(id, orgID)
}

func (s *MaintenanceService) List(orgID uint, page, pageSize int) ([]repository.MaintenancePlan, int, error) {
	return s.repo.ListMaintenancePlans(orgID, page, pageSize)
}

func (s *MaintenanceService) Update(plan *repository.MaintenancePlan) error {
	return s.repo.UpdateMaintenancePlan(plan)
}

func (s *MaintenanceService) Delete(id, orgID uint) error {
	return s.repo.DeleteMaintenancePlan(id, orgID)
}
