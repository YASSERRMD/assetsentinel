package services

import (
	"errors"
	"time"

	"assetsentinel/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewAuthService(repo *repository.Repository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

type Claims struct {
	UserID         uint   `json:"user_id"`
	OrganizationID uint   `json:"organization_id"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(email, password, fullName, role string, orgID uint) (*repository.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		OrganizationID: orgID,
		Email:          email,
		PasswordHash:   string(hash),
		FullName:       fullName,
		Role:           role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *repository.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := &Claims{
		UserID:         user.ID,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

type AssetService struct {
	repo *repository.Repository
}

func NewAssetService(repo *repository.Repository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) Create(asset *repository.Asset) error {
	return s.repo.CreateAsset(asset)
}

func (s *AssetService) Get(id, orgID uint) (*repository.Asset, error) {
	return s.repo.GetAsset(id, orgID)
}

func (s *AssetService) List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error) {
	return s.repo.ListAssets(orgID, page, pageSize, status, category)
}

func (s *AssetService) Update(asset *repository.Asset) error {
	return s.repo.UpdateAsset(asset)
}

func (s *AssetService) Delete(id, orgID uint) error {
	return s.repo.SoftDeleteAsset(id, orgID)
}

type MaintenanceService struct {
	repo *repository.Repository
	hub  interface {
		BroadcastToOrg(orgID uint, message map[string]interface{})
	}
}

func NewMaintenanceService(repo *repository.Repository, hub interface {
	BroadcastToOrg(orgID uint, message map[string]interface{})
}) *MaintenanceService {
	return &MaintenanceService{repo: repo, hub: hub}
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

type WorkOrderService struct {
	repo *repository.Repository
	hub  interface {
		BroadcastToOrg(orgID uint, message map[string]interface{})
	}
}

func NewWorkOrderService(repo *repository.Repository, hub interface {
	BroadcastToOrg(orgID uint, message map[string]interface{})
}) *WorkOrderService {
	return &WorkOrderService{repo: repo, hub: hub}
}

func (s *WorkOrderService) Create(wo *repository.WorkOrder) error {
	if err := s.repo.CreateWorkOrder(wo); err != nil {
		return err
	}
	if s.hub != nil {
		s.hub.BroadcastToOrg(wo.OrganizationID, map[string]interface{}{
			"type":       "work_order_created",
			"work_order": wo,
		})
	}
	return nil
}

func (s *WorkOrderService) Get(id, orgID uint) (*repository.WorkOrder, error) {
	return s.repo.GetWorkOrder(id, orgID)
}

func (s *WorkOrderService) List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error) {
	return s.repo.ListWorkOrders(orgID, page, pageSize, status)
}

func (s *WorkOrderService) Update(wo *repository.WorkOrder, oldStatus string) error {
	if err := s.repo.UpdateWorkOrder(wo); err != nil {
		return err
	}
	if s.hub != nil && oldStatus != wo.Status {
		s.hub.BroadcastToOrg(wo.OrganizationID, map[string]interface{}{
			"type":       "work_order_status_change",
			"work_order": wo,
			"old_status": oldStatus,
			"new_status": wo.Status,
		})
	}
	return nil
}

func (s *WorkOrderService) Delete(id, orgID uint) error {
	return s.repo.DeleteWorkOrder(id, orgID)
}

type InventoryService struct {
	repo *repository.Repository
	hub  interface {
		BroadcastToOrg(orgID uint, message map[string]interface{})
	}
}

func NewInventoryService(repo *repository.Repository, hub interface {
	BroadcastToOrg(orgID uint, message map[string]interface{})
}) *InventoryService {
	return &InventoryService{repo: repo, hub: hub}
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
	oldPart, err := s.repo.GetInventoryPart(part.ID, part.OrganizationID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateInventoryPart(part); err != nil {
		return err
	}

	if s.hub != nil && oldPart.Quantity > part.MinThreshold && part.Quantity <= part.MinThreshold {
		s.hub.BroadcastToOrg(part.OrganizationID, map[string]interface{}{
			"type": "low_inventory",
			"part": part,
		})
	}
	return nil
}

func (s *InventoryService) Deduct(partID, orgID, quantity int) (int, error) {
	return s.repo.DeductInventory(partID, orgID, quantity)
}

func (s *InventoryService) Delete(id, orgID uint) error {
	return s.repo.DeleteInventoryPart(id, orgID)
}

type DepreciationService struct {
	repo *repository.Repository
}

func NewDepreciationService(repo *repository.Repository) *DepreciationService {
	return &DepreciationService{repo: repo}
}

func (s *DepreciationService) CalculateStraightLine(asset *repository.Asset, usefulLifeYears int) []repository.AssetDepreciation {
	annualDepreciation := asset.PurchaseCost / float64(usefulLifeYears)
	year := time.Now().Year()

	var deprs []repository.AssetDepreciation
	for i := 0; i < usefulLifeYears; i++ {
		deprs = append(deprs, repository.AssetDepreciation{
			OrganizationID:     asset.OrganizationID,
			AssetID:            asset.ID,
			Year:               year + i,
			DepreciationAmount: annualDepreciation,
		})
	}
	return deprs
}

func (s *DepreciationService) SaveDepreciation(depr *repository.AssetDepreciation) error {
	return s.repo.CreateDepreciation(depr)
}

func (s *DepreciationService) GetAssetDepreciation(assetID, orgID uint) ([]repository.AssetDepreciation, error) {
	return s.repo.GetAssetDepreciation(assetID, orgID)
}

func (s *DepreciationService) GetAssetCosts(assetID, orgID uint) (float64, error) {
	return s.repo.GetAssetMaintenanceCost(assetID, orgID)
}

func (s *DepreciationService) GetAllCosts(orgID uint) ([]struct {
	AssetID   uint
	AssetName string
	TotalCost float64
}, error) {
	return s.repo.GetAllAssetCosts(orgID)
}
