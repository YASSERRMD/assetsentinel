package repository

import (
	"time"
)

type Organization struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID             uint      `json:"id"`
	OrganizationID uint      `json:"organization_id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	FullName       string    `json:"full_name"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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

type MaintenancePlan struct {
	ID                     uint       `json:"id"`
	OrganizationID         uint       `json:"organization_id"`
	AssetID                uint       `json:"asset_id"`
	FrequencyDays          int        `json:"frequency_days"`
	EstimatedDurationHours *float64   `json:"estimated_duration_hours"`
	AssignedRole           *string    `json:"assigned_role"`
	LastMaintenanceDate    *time.Time `json:"last_maintenance_date"`
	NextMaintenanceDate    time.Time  `json:"next_maintenance_date"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

type MaintenanceTask struct {
	ID                uint       `json:"id"`
	OrganizationID    uint       `json:"organization_id"`
	MaintenancePlanID uint       `json:"maintenance_plan_id"`
	AssetID           uint       `json:"asset_id"`
	ScheduledDate     time.Time  `json:"scheduled_date"`
	Status            string     `json:"status"`
	CompletedDate     *time.Time `json:"completed_date"`
	Notes             *string    `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type WorkOrder struct {
	ID             uint       `json:"id"`
	OrganizationID uint       `json:"organization_id"`
	AssetID        uint       `json:"asset_id"`
	TechnicianID   *uint      `json:"technician_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
	TotalCost      float64    `json:"total_cost"`
	Notes          *string    `json:"notes"`
	CreatedBy      *uint      `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type WorkOrderPart struct {
	ID          uint      `json:"id"`
	WorkOrderID uint      `json:"work_order_id"`
	PartID      uint      `json:"part_id"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type InventoryPart struct {
	ID             uint       `json:"id"`
	OrganizationID uint       `json:"organization_id"`
	Name           string     `json:"name"`
	SKU            string     `json:"sku"`
	Quantity       int        `json:"quantity"`
	MinThreshold   int        `json:"min_threshold"`
	CostPerUnit    float64    `json:"cost_per_unit"`
	Location       *string    `json:"location"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AssetDepreciation struct {
	ID                 uint      `json:"id"`
	OrganizationID     uint      `json:"organization_id"`
	AssetID            uint      `json:"asset_id"`
	Year               int       `json:"year"`
	DepreciationAmount float64   `json:"depreciation_amount"`
	CreatedAt          time.Time `json:"created_at"`
}

type AuditLog struct {
	ID             uint      `json:"id"`
	OrganizationID uint      `json:"organization_id"`
	UserID         *uint     `json:"user_id"`
	TableName      string    `json:"table_name"`
	RecordID       uint      `json:"record_id"`
	Action         string    `json:"action"`
	OldValues      *string   `json:"old_values"`
	NewValues      *string   `json:"new_values"`
	CreatedAt      time.Time `json:"created_at"`
}

func (db *DB) CreateOrganization(org *Organization) error {
	result, err := db.Exec(`INSERT INTO organizations (name) VALUES (?)`, org.Name)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	org.ID = uint(id)
	return nil
}

func (db *DB) GetOrganization(id uint) (*Organization, error) {
	org := &Organization{}
	err := db.QueryRow(`SELECT id, name, created_at, updated_at FROM organizations WHERE id = ?`, id).
		Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	return org, err
}

func (db *DB) ListOrganizations() ([]Organization, error) {
	rows, err := db.Query(`SELECT id, name, created_at, updated_at FROM organizations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (db *DB) UpdateOrganization(org *Organization) error {
	_, err := db.Exec(`UPDATE organizations SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, org.Name, org.ID)
	return err
}

func (db *DB) CreateUser(user *User) error {
	result, err := db.Exec(`INSERT INTO users (organization_id, email, password_hash, full_name, role) VALUES (?, ?, ?, ?, ?)`,
		user.OrganizationID, user.Email, user.PasswordHash, user.FullName, user.Role)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	user.ID = uint(id)
	return nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.QueryRow(`SELECT id, organization_id, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE email = ?`, email).
		Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (db *DB) GetUser(id uint) (*User, error) {
	user := &User{}
	err := db.QueryRow(`SELECT id, organization_id, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.OrganizationID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (db *DB) ListUsers(orgID uint) ([]User, error) {
	rows, err := db.Query(`SELECT id, organization_id, email, full_name, role, created_at, updated_at FROM users WHERE organization_id = ?`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.OrganizationID, &user.Email, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) UpdateUser(user *User) error {
	_, err := db.Exec(`UPDATE users SET email = ?, full_name = ?, role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		user.Email, user.FullName, user.Role, user.ID)
	return err
}

func (db *DB) DeleteUser(id uint) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
