package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	*DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{db}
}

type Asset struct {
	ID               uint       `json:"id"`
	OrganizationID   uint       `json:"organization_id"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	SerialNumber     *string    `json:"serial_number"`
	InstallationDate *time.Time `json:"installation_date"`
	Location         *string    `json:"location"`
	PurchaseCost     float64    `json:"purchase_cost"`
	WarrantyExpiry   *time.Time `json:"warranty_expiry"`
	Status           string     `json:"status"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (r *Repository) CreateAsset(asset *Asset) error {
	result, err := r.Exec(`INSERT INTO assets (organization_id, name, category, serial_number, installation_date, location, purchase_cost, warranty_expiry, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		asset.OrganizationID, asset.Name, asset.Category, asset.SerialNumber, asset.InstallationDate, asset.Location, asset.PurchaseCost, asset.WarrantyExpiry, asset.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	asset.ID = uint(id)
	return nil
}

func (r *Repository) GetAsset(id, orgID uint) (*Asset, error) {
	asset := &Asset{}
	err := r.QueryRow(`SELECT id, organization_id, name, category, serial_number, installation_date, location, purchase_cost, warranty_expiry, status, deleted_at, created_at, updated_at 
		FROM assets WHERE id = ? AND organization_id = ? AND deleted_at IS NULL`, id, orgID).
		Scan(&asset.ID, &asset.OrganizationID, &asset.Name, &asset.Category, &asset.SerialNumber, &asset.InstallationDate, &asset.Location,
			&asset.PurchaseCost, &asset.WarrantyExpiry, &asset.Status, &asset.DeletedAt, &asset.CreatedAt, &asset.UpdatedAt)
	return asset, err
}

func (r *Repository) ListAssets(orgID uint, page, pageSize int, status, category string) ([]Asset, int, error) {
	offset := (page - 1) * pageSize

	var count int
	countQuery := `SELECT COUNT(*) FROM assets WHERE organization_id = ? AND deleted_at IS NULL`
	queryArgs := []interface{}{orgID}

	if status != "" {
		countQuery += ` AND status = ?`
		queryArgs = append(queryArgs, status)
	}
	if category != "" {
		countQuery += ` AND category = ?`
		queryArgs = append(queryArgs, category)
	}

	err := r.QueryRow(countQuery, queryArgs...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, organization_id, name, category, serial_number, installation_date, location, purchase_cost, warranty_expiry, status, created_at, updated_at 
		FROM assets WHERE organization_id = ? AND deleted_at IS NULL`
	queryArgs = []interface{}{orgID}

	if status != "" {
		query += ` AND status = ?`
		queryArgs = append(queryArgs, status)
	}
	if category != "" {
		query += ` AND category = ?`
		queryArgs = append(queryArgs, category)
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := r.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		if err := rows.Scan(&asset.ID, &asset.OrganizationID, &asset.Name, &asset.Category, &asset.SerialNumber,
			&asset.InstallationDate, &asset.Location, &asset.PurchaseCost, &asset.WarrantyExpiry, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, 0, err
		}
		assets = append(assets, asset)
	}
	return assets, count, nil
}

func (r *Repository) UpdateAsset(asset *Asset) error {
	_, err := r.Exec(`UPDATE assets SET name = ?, category = ?, serial_number = ?, installation_date = ?, location = ?, purchase_cost = ?, warranty_expiry = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`,
		asset.Name, asset.Category, asset.SerialNumber, asset.InstallationDate, asset.Location, asset.PurchaseCost, asset.WarrantyExpiry, asset.Status, asset.ID, asset.OrganizationID)
	return err
}

func (r *Repository) SoftDeleteAsset(id, orgID uint) error {
	_, err := r.Exec(`UPDATE assets SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`, id, orgID)
	return err
}

func (r *Repository) GetDashboardStats(orgID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var assetCount int
	err := r.QueryRow(`SELECT COUNT(*) FROM assets WHERE organization_id = ? AND deleted_at IS NULL`, orgID).Scan(&assetCount)
	if err != nil {
		return nil, err
	}
	stats["asset_count"] = assetCount

	var lowStockCount int
	err = r.QueryRow(`SELECT COUNT(*) FROM inventory_parts WHERE organization_id = ? AND quantity <= min_threshold AND deleted_at IS NULL`, orgID).Scan(&lowStockCount)
	if err != nil {
		return nil, err
	}
	stats["low_stock"] = lowStockCount

	var totalCost float64
	err = r.QueryRow(`SELECT COALESCE(SUM(total_cost), 0) FROM work_orders WHERE organization_id = ?`, orgID).Scan(&totalCost)
	if err != nil {
		return nil, err
	}
	stats["total_costs"] = totalCost

	return stats, nil
}
