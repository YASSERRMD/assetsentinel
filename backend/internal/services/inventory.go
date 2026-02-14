package services

import (
	"assetsentinel/internal/repository"
)

type InventoryService struct {
	repo *repository.Repository
}

func NewInventoryService(repo *repository.Repository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) Create(part *repository.InventoryPart) error {
	return s.repo.CreateInventoryPart(part)
}

func (s *InventoryService) Get(id, orgID uint) (*repository.InventoryPart, error) {
	return s.repo.GetInventoryPart(id, orgID)
}

func (s *InventoryService) List(orgID uint, page, pageSize int) ([]repository.InventoryPart, int, error) {
	return s.repo.ListInventoryParts(orgID, page, pageSize)
}

func (s *InventoryService) GetLowStock(orgID uint) ([]repository.InventoryPart, error) {
	return s.repo.GetLowStockParts(orgID)
}

func (s *InventoryService) Update(part *repository.InventoryPart) error {
	return s.repo.UpdateInventoryPart(part)
}

func (s *InventoryService) Delete(id, orgID uint) error {
	return s.repo.DeleteInventoryPart(id, orgID)
}

type WorkOrderService struct {
	repo *repository.Repository
}

func NewWorkOrderService(repo *repository.Repository) *WorkOrderService {
	return &WorkOrderService{repo: repo}
}

func (s *WorkOrderService) Create(wo *repository.WorkOrder) error {
	return s.repo.CreateWorkOrder(wo)
}

func (s *WorkOrderService) Get(id, orgID uint) (*repository.WorkOrder, error) {
	return s.repo.GetWorkOrder(id, orgID)
}

func (s *WorkOrderService) List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error) {
	return s.repo.ListWorkOrders(orgID, page, pageSize, status)
}

func (s *WorkOrderService) Update(wo *repository.WorkOrder) error {
	return s.repo.UpdateWorkOrder(wo)
}

func (s *WorkOrderService) Delete(id, orgID uint) error {
	return s.repo.DeleteWorkOrder(id, orgID)
}
