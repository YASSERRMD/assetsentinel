package handlers

import (
	"net/http"
	"strconv"

	"assetsentinel/internal/middleware"
	"assetsentinel/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService interface {
		Login(email, password string) (string, *repository.User, error)
		Register(email, password, fullName, role string, orgID uint) (*repository.User, error)
	}
}

func NewAuthHandler(authService interface {
	Login(email, password string) (string, *repository.User, error)
	Register(email, password, fullName, role string, orgID uint) (*repository.User, error)
}) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"full_name":       user.FullName,
			"role":            user.Role,
			"organization_id": user.OrganizationID,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Role     string `json:"role" binding:"required"`
		OrgID    uint   `json:"organization_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.FullName, req.Role, req.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              user.ID,
		"email":           user.Email,
		"full_name":       user.FullName,
		"role":            user.Role,
		"organization_id": user.OrganizationID,
	})
}

type AssetHandler struct {
	assetService interface {
		Create(asset *repository.Asset) error
		Get(id, orgID uint) (*repository.Asset, error)
		List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error)
		Update(asset *repository.Asset) error
		Delete(id, orgID uint) error
	}
}

func NewAssetHandler(assetService interface {
	Create(asset *repository.Asset) error
	Get(id, orgID uint) (*repository.Asset, error)
	List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error)
	Update(asset *repository.Asset) error
	Delete(id, orgID uint) error
}) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) Create(c *gin.Context) {
	var asset repository.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.OrganizationID = middleware.GetOrganizationID(c)

	if err := h.assetService.Create(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	asset, err := h.assetService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	category := c.Query("category")

	assets, total, err := h.assetService.List(orgID, page, pageSize, status, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      assets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *AssetHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	var asset repository.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = uint(id)
	asset.OrganizationID = orgID

	if err := h.assetService.Update(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.assetService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted"})
}

type MaintenanceHandler struct {
	maintenanceService interface {
		Create(plan *repository.MaintenancePlan) error
		Get(id, orgID uint) (*repository.MaintenancePlan, error)
		List(orgID uint, page, pageSize int) ([]repository.MaintenancePlan, int, error)
		Update(plan *repository.MaintenancePlan) error
		Delete(id, orgID uint) error
	}
}

func NewMaintenanceHandler(maintenanceService interface {
	Create(plan *repository.MaintenancePlan) error
	Get(id, orgID uint) (*repository.MaintenancePlan, error)
	List(orgID uint, page, pageSize int) ([]repository.MaintenancePlan, int, error)
	Update(plan *repository.MaintenancePlan) error
	Delete(id, orgID uint) error
}) *MaintenanceHandler {
	return &MaintenanceHandler{maintenanceService: maintenanceService}
}

func (h *MaintenanceHandler) Create(c *gin.Context) {
	var plan repository.MaintenancePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.OrganizationID = middleware.GetOrganizationID(c)

	if err := h.maintenanceService.Create(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *MaintenanceHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	plan, err := h.maintenanceService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *MaintenanceHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	plans, total, err := h.maintenanceService.List(orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      plans,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *MaintenanceHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	var plan repository.MaintenancePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.ID = uint(id)
	plan.OrganizationID = orgID

	if err := h.maintenanceService.Update(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *MaintenanceHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.maintenanceService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance plan deleted"})
}

type WorkOrderHandler struct {
	workOrderService interface {
		Create(wo *repository.WorkOrder) error
		Get(id, orgID uint) (*repository.WorkOrder, error)
		List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error)
		Update(wo *repository.WorkOrder, oldStatus string) error
		Delete(id, orgID uint) error
	}
}

func NewWorkOrderHandler(workOrderService interface {
	Create(wo *repository.WorkOrder) error
	Get(id, orgID uint) (*repository.WorkOrder, error)
	List(orgID uint, page, pageSize int, status string) ([]repository.WorkOrder, int, error)
	Update(wo *repository.WorkOrder, oldStatus string) error
	Delete(id, orgID uint) error
}) *WorkOrderHandler {
	return &WorkOrderHandler{workOrderService: workOrderService}
}

func (h *WorkOrderHandler) Create(c *gin.Context) {
	var wo repository.WorkOrder
	if err := c.ShouldBindJSON(&wo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wo.OrganizationID = middleware.GetOrganizationID(c)
	wo.CreatedBy = func() *uint { id := middleware.GetUserID(c); return &id }()

	if err := h.workOrderService.Create(&wo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wo)
}

func (h *WorkOrderHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	wo, err := h.workOrderService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work order not found"})
		return
	}

	c.JSON(http.StatusOK, wo)
}

func (h *WorkOrderHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	orders, total, err := h.workOrderService.List(orgID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *WorkOrderHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	existing, _ := h.workOrderService.Get(uint(id), orgID)
	oldStatus := ""
	if existing != nil {
		oldStatus = existing.Status
	}

	var wo repository.WorkOrder
	if err := c.ShouldBindJSON(&wo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wo.ID = uint(id)
	wo.OrganizationID = orgID

	if err := h.workOrderService.Update(&wo, oldStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wo)
}

func (h *WorkOrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.workOrderService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work order deleted"})
}

type InventoryHandler struct {
	inventoryService interface {
		Create(part *repository.InventoryPart) error
		Get(id, orgID uint) (*repository.InventoryPart, error)
		List(orgID uint, page, pageSize int) ([]repository.InventoryPart, int, error)
		GetLowStock(orgID uint) ([]repository.InventoryPart, error)
		Update(part *repository.InventoryPart) error
		Deduct(partID, orgID, quantity int) (int, error)
		Delete(id, orgID uint) error
	}
}

func NewInventoryHandler(inventoryService interface {
	Create(part *repository.InventoryPart) error
	Get(id, orgID uint) (*repository.InventoryPart, error)
	List(orgID uint, page, pageSize int) ([]repository.InventoryPart, int, error)
	GetLowStock(orgID uint) ([]repository.InventoryPart, error)
	Update(part *repository.InventoryPart) error
	Deduct(partID, orgID, quantity int) (int, error)
	Delete(id, orgID uint) error
}) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) Create(c *gin.Context) {
	var part repository.InventoryPart
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	part.OrganizationID = middleware.GetOrganizationID(c)

	if err := h.inventoryService.Create(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, part)
}

func (h *InventoryHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	part, err := h.inventoryService.Get(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory part not found"})
		return
	}

	c.JSON(http.StatusOK, part)
}

func (h *InventoryHandler) List(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	parts, total, err := h.inventoryService.List(orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      parts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *InventoryHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	var part repository.InventoryPart
	if err := c.ShouldBindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	part.ID = uint(id)
	part.OrganizationID = orgID

	if err := h.inventoryService.Update(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, part)
}

func (h *InventoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	if err := h.inventoryService.Delete(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory part deleted"})
}

type DepreciationHandler struct {
	depreciationService interface {
		CalculateStraightLine(asset *repository.Asset, usefulLifeYears int) []repository.AssetDepreciation
		SaveDepreciation(depr *repository.AssetDepreciation) error
		GetAssetDepreciation(assetID, orgID uint) ([]repository.AssetDepreciation, error)
		GetAssetCosts(assetID, orgID uint) (float64, error)
		GetAllCosts(orgID uint) ([]struct {
			AssetID   uint
			AssetName string
			TotalCost float64
		}, error)
	}
}

func NewDepreciationHandler(depreciationService interface {
	CalculateStraightLine(asset *repository.Asset, usefulLifeYears int) []repository.AssetDepreciation
	SaveDepreciation(depr *repository.AssetDepreciation) error
	GetAssetDepreciation(assetID, orgID uint) ([]repository.AssetDepreciation, error)
	GetAssetCosts(assetID, orgID uint) (float64, error)
	GetAllCosts(orgID uint) ([]struct {
		AssetID   uint
		AssetName string
		TotalCost float64
	}, error)
}) *DepreciationHandler {
	return &DepreciationHandler{depreciationService: depreciationService}
}

func (h *DepreciationHandler) GetAssetDepreciation(c *gin.Context) {
	assetID, _ := strconv.ParseUint(c.Param("asset_id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	deprs, err := h.depreciationService.GetAssetDepreciation(uint(assetID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deprs)
}

func (h *DepreciationHandler) GetAssetCosts(c *gin.Context) {
	assetID, _ := strconv.ParseUint(c.Param("asset_id"), 10, 32)
	orgID := middleware.GetOrganizationID(c)

	costs, err := h.depreciationService.GetAssetCosts(uint(assetID), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_maintenance_cost": costs})
}

func (h *DepreciationHandler) GetAllCosts(c *gin.Context) {
	orgID := middleware.GetOrganizationID(c)

	costs, err := h.depreciationService.GetAllCosts(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, costs)
}

func GetDashboard(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := middleware.GetOrganizationID(c)

		stats, err := repo.GetDashboardStats(orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

func GetAuditLogs(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := middleware.GetOrganizationID(c)
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		tableName := c.Query("table")

		logs, total, err := repo.GetAuditLogs(orgID, page, pageSize, tableName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

func ListOrganizations(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs, err := repo.ListOrganizations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orgs)
	}
}

func CreateOrganization(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org repository.Organization
		if err := c.ShouldBindJSON(&org); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.CreateOrganization(&org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, org)
	}
}

func GetOrganization(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		org, err := repo.GetOrganization(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func UpdateOrganization(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		var org repository.Organization
		if err := c.ShouldBindJSON(&org); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org.ID = uint(id)

		if err := repo.UpdateOrganization(&org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func ListUsers(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := middleware.GetOrganizationID(c)

		users, err := repo.ListUsers(orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func CreateUser(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := middleware.GetOrganizationID(c)

		var user repository.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.OrganizationID = orgID

		if err := repo.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func GetUser(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		user, err := repo.GetUser(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		var user repository.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.ID = uint(id)

		if err := repo.UpdateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(repo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		if err := repo.DeleteUser(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}
