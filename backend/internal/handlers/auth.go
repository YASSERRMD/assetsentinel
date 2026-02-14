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
