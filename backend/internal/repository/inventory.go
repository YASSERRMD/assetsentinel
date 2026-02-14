package repository

type InventoryPart struct {
	ID             uint    `json:"id"`
	OrganizationID uint    `json:"organization_id"`
	Name           string  `json:"name"`
	SKU            string  `json:"sku"`
	Quantity       int     `json:"quantity"`
	MinThreshold   int     `json:"min_threshold"`
	CostPerUnit    float64 `json:"cost_per_unit"`
	Location       *string `json:"location"`
	DeletedAt      *string `json:"deleted_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type WorkOrder struct {
	ID             uint    `json:"id"`
	OrganizationID uint    `json:"organization_id"`
	AssetID        uint    `json:"asset_id"`
	TechnicianID   *uint   `json:"technician_id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	Status         string  `json:"status"`
	Priority       string  `json:"priority"`
	ScheduledStart *string `json:"scheduled_start"`
	ScheduledEnd   *string `json:"scheduled_end"`
	ActualStart    *string `json:"actual_start"`
	ActualEnd      *string `json:"actual_end"`
	TotalCost      float64 `json:"total_cost"`
	Notes          *string `json:"notes"`
	CreatedBy      *uint   `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func (r *Repository) CreateInventoryPart(part *InventoryPart) error {
	result, err := r.Exec(`INSERT INTO inventory_parts (organization_id, name, sku, quantity, min_threshold, cost_per_unit, location) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		part.OrganizationID, part.Name, part.SKU, part.Quantity, part.MinThreshold, part.CostPerUnit, part.Location)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	part.ID = uint(id)
	return nil
}

func (r *Repository) GetInventoryPart(id, orgID uint) (*InventoryPart, error) {
	part := &InventoryPart{}
	err := r.QueryRow(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, deleted_at, created_at, updated_at FROM inventory_parts WHERE id = ? AND organization_id = ? AND deleted_at IS NULL`, id, orgID).
		Scan(&part.ID, &part.OrganizationID, &part.Name, &part.SKU, &part.Quantity, &part.MinThreshold, &part.CostPerUnit, &part.Location, &part.DeletedAt, &part.CreatedAt, &part.UpdatedAt)
	return part, err
}

func (r *Repository) ListInventoryParts(orgID uint, page, pageSize int) ([]InventoryPart, int, error) {
	offset := (page - 1) * pageSize
	var count int
	r.QueryRow(`SELECT COUNT(*) FROM inventory_parts WHERE organization_id = ? AND deleted_at IS NULL`, orgID).Scan(&count)

	rows, err := r.Query(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, created_at, updated_at FROM inventory_parts WHERE organization_id = ? AND deleted_at IS NULL ORDER BY name ASC LIMIT ? OFFSET ?`, orgID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var parts []InventoryPart
	for rows.Next() {
		var part InventoryPart
		if err := rows.Scan(&part.ID, &part.OrganizationID, &part.Name, &part.SKU, &part.Quantity, &part.MinThreshold, &part.CostPerUnit, &part.Location, &part.CreatedAt, &part.UpdatedAt); err != nil {
			return nil, 0, err
		}
		parts = append(parts, part)
	}
	return parts, count, nil
}

func (r *Repository) UpdateInventoryPart(part *InventoryPart) error {
	_, err := r.Exec(`UPDATE inventory_parts SET name = ?, sku = ?, quantity = ?, min_threshold = ?, cost_per_unit = ?, location = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`,
		part.Name, part.SKU, part.Quantity, part.MinThreshold, part.CostPerUnit, part.Location, part.ID, part.OrganizationID)
	return err
}

func (r *Repository) DeleteInventoryPart(id, orgID uint) error {
	_, err := r.Exec(`UPDATE inventory_parts SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`, id, orgID)
	return err
}

func (r *Repository) GetLowStockParts(orgID uint) ([]InventoryPart, error) {
	rows, err := r.Query(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, created_at, updated_at FROM inventory_parts WHERE organization_id = ? AND quantity <= min_threshold AND deleted_at IS NULL`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []InventoryPart
	for rows.Next() {
		var part InventoryPart
		if err := rows.Scan(&part.ID, &part.OrganizationID, &part.Name, &part.SKU, &part.Quantity, &part.MinThreshold, &part.CostPerUnit, &part.Location, &part.CreatedAt, &part.UpdatedAt); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
}

func (r *Repository) CreateWorkOrder(wo *WorkOrder) error {
	result, err := r.Exec(`INSERT INTO work_orders (organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, notes, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		wo.OrganizationID, wo.AssetID, wo.TechnicianID, wo.Title, wo.Description, wo.Status, wo.Priority, wo.ScheduledStart, wo.ScheduledEnd, wo.Notes, wo.CreatedBy)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	wo.ID = uint(id)
	return nil
}

func (r *Repository) GetWorkOrder(id, orgID uint) (*WorkOrder, error) {
	wo := &WorkOrder{}
	err := r.QueryRow(`SELECT id, organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, actual_start, actual_end, total_cost, notes, created_by, created_at, updated_at FROM work_orders WHERE id = ? AND organization_id = ?`, id, orgID).
		Scan(&wo.ID, &wo.OrganizationID, &wo.AssetID, &wo.TechnicianID, &wo.Title, &wo.Description, &wo.Status, &wo.Priority, &wo.ScheduledStart, &wo.ScheduledEnd, &wo.ActualStart, &wo.ActualEnd, &wo.TotalCost, &wo.Notes, &wo.CreatedBy, &wo.CreatedAt, &wo.UpdatedAt)
	return wo, err
}

func (r *Repository) ListWorkOrders(orgID uint, page, pageSize int, status string) ([]WorkOrder, int, error) {
	offset := (page - 1) * pageSize
	var count int
	query := `SELECT COUNT(*) FROM work_orders WHERE organization_id = ?`
	args := []interface{}{orgID}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	r.QueryRow(query, args...).Scan(&count)

	query = `SELECT id, organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, actual_start, actual_end, total_cost, notes, created_by, created_at, updated_at FROM work_orders WHERE organization_id = ?`
	args = []interface{}{orgID}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []WorkOrder
	for rows.Next() {
		var wo WorkOrder
		if err := rows.Scan(&wo.ID, &wo.OrganizationID, &wo.AssetID, &wo.TechnicianID, &wo.Title, &wo.Description, &wo.Status, &wo.Priority, &wo.ScheduledStart, &wo.ScheduledEnd, &wo.ActualStart, &wo.ActualEnd, &wo.TotalCost, &wo.Notes, &wo.CreatedBy, &wo.CreatedAt, &wo.UpdatedAt); err != nil {
			return nil, 0, err
		}
		orders = append(orders, wo)
	}
	return orders, count, nil
}

func (r *Repository) UpdateWorkOrder(wo *WorkOrder) error {
	_, err := r.Exec(`UPDATE work_orders SET asset_id = ?, technician_id = ?, title = ?, description = ?, status = ?, priority = ?, scheduled_start = ?, scheduled_end = ?, actual_start = ?, actual_end = ?, total_cost = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`,
		wo.AssetID, wo.TechnicianID, wo.Title, wo.Description, wo.Status, wo.Priority, wo.ScheduledStart, wo.ScheduledEnd, wo.ActualStart, wo.ActualEnd, wo.TotalCost, wo.Notes, wo.ID, wo.OrganizationID)
	return err
}

func (r *Repository) DeleteWorkOrder(id, orgID uint) error {
	_, err := r.Exec(`DELETE FROM work_orders WHERE id = ? AND organization_id = ?`, id, orgID)
	return err
}

func (r *Repository) GetAllAssetCosts(orgID uint) ([]struct {
	AssetID   uint
	AssetName string
	TotalCost float64
}, error) {
	rows, err := r.Query(`SELECT a.id, a.name, COALESCE(SUM(wo.total_cost), 0) as total_cost FROM assets a LEFT JOIN work_orders wo ON a.id = wo.asset_id WHERE a.organization_id = ? AND a.deleted_at IS NULL GROUP BY a.id, a.name`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		AssetID   uint
		AssetName string
		TotalCost float64
	}
	for rows.Next() {
		var r struct {
			AssetID   uint
			AssetName string
			TotalCost float64
		}
		if err := rows.Scan(&r.AssetID, &r.AssetName, &r.TotalCost); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (r *Repository) GetAssetMaintenanceCost(assetID, orgID uint) (float64, error) {
	var totalCost float64
	err := r.QueryRow(`SELECT COALESCE(SUM(total_cost), 0) FROM work_orders WHERE asset_id = ? AND organization_id = ?`, assetID, orgID).Scan(&totalCost)
	return totalCost, err
}
