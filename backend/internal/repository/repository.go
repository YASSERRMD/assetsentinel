package repository

import (
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

func (r *Repository) CreateMaintenancePlan(plan *MaintenancePlan) error {
	result, err := r.Exec(`INSERT INTO maintenance_plans (organization_id, asset_id, frequency_days, estimated_duration_hours, assigned_role, last_maintenance_date, next_maintenance_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		plan.OrganizationID, plan.AssetID, plan.FrequencyDays, plan.EstimatedDurationHours, plan.AssignedRole, plan.LastMaintenanceDate, plan.NextMaintenanceDate)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	plan.ID = uint(id)
	return nil
}

func (r *Repository) GetMaintenancePlan(id, orgID uint) (*MaintenancePlan, error) {
	plan := &MaintenancePlan{}
	err := r.QueryRow(`SELECT id, organization_id, asset_id, frequency_days, estimated_duration_hours, assigned_role, last_maintenance_date, next_maintenance_date, created_at, updated_at 
		FROM maintenance_plans WHERE id = ? AND organization_id = ?`, id, orgID).
		Scan(&plan.ID, &plan.OrganizationID, &plan.AssetID, &plan.FrequencyDays, &plan.EstimatedDurationHours, &plan.AssignedRole, &plan.LastMaintenanceDate, &plan.NextMaintenanceDate, &plan.CreatedAt, &plan.UpdatedAt)
	return plan, err
}

func (r *Repository) ListMaintenancePlans(orgID uint, page, pageSize int) ([]MaintenancePlan, int, error) {
	offset := (page - 1) * pageSize

	var count int
	err := r.QueryRow(`SELECT COUNT(*) FROM maintenance_plans WHERE organization_id = ?`, orgID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.Query(`SELECT id, organization_id, asset_id, frequency_days, estimated_duration_hours, assigned_role, last_maintenance_date, next_maintenance_date, created_at, updated_at 
		FROM maintenance_plans WHERE organization_id = ? ORDER BY next_maintenance_date ASC LIMIT ? OFFSET ?`, orgID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var plans []MaintenancePlan
	for rows.Next() {
		var plan MaintenancePlan
		if err := rows.Scan(&plan.ID, &plan.OrganizationID, &plan.AssetID, &plan.FrequencyDays, &plan.EstimatedDurationHours, &plan.AssignedRole, &plan.LastMaintenanceDate, &plan.NextMaintenanceDate, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, 0, err
		}
		plans = append(plans, plan)
	}
	return plans, count, nil
}

func (r *Repository) UpdateMaintenancePlan(plan *MaintenancePlan) error {
	_, err := r.Exec(`UPDATE maintenance_plans SET asset_id = ?, frequency_days = ?, estimated_duration_hours = ?, assigned_role = ?, last_maintenance_date = ?, next_maintenance_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`,
		plan.AssetID, plan.FrequencyDays, plan.EstimatedDurationHours, plan.AssignedRole, plan.LastMaintenanceDate, plan.NextMaintenanceDate, plan.ID, plan.OrganizationID)
	return err
}

func (r *Repository) DeleteMaintenancePlan(id, orgID uint) error {
	_, err := r.Exec(`DELETE FROM maintenance_plans WHERE id = ? AND organization_id = ?`, id, orgID)
	return err
}

func (r *Repository) GetUpcomingMaintenanceTasks(orgID uint, daysAhead int) ([]MaintenanceTask, error) {
	futureDate := time.Now().AddDate(0, 0, daysAhead)
	rows, err := r.Query(`SELECT mt.id, mt.organization_id, mt.maintenance_plan_id, mt.asset_id, mt.scheduled_date, mt.status, mt.completed_date, mt.notes, mt.created_at, mt.updated_at 
		FROM maintenance_tasks mt 
		JOIN maintenance_plans mp ON mt.maintenance_plan_id = mp.id 
		WHERE mp.organization_id = ? AND mt.scheduled_date <= ? AND mt.status = 'pending'`, orgID, futureDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []MaintenanceTask
	for rows.Next() {
		var task MaintenanceTask
		if err := rows.Scan(&task.ID, &task.OrganizationID, &task.MaintenancePlanID, &task.AssetID, &task.ScheduledDate, &task.Status, &task.CompletedDate, &task.Notes, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetOverdueMaintenanceTasks(orgID uint) ([]MaintenanceTask, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := r.Query(`SELECT mt.id, mt.organization_id, mt.maintenance_plan_id, mt.asset_id, mt.scheduled_date, mt.status, mt.completed_date, mt.notes, mt.created_at, mt.updated_at 
		FROM maintenance_tasks mt 
		JOIN maintenance_plans mp ON mt.maintenance_plan_id = mp.id 
		WHERE mp.organization_id = ? AND mt.scheduled_date < ? AND mt.status IN ('pending', 'in_progress')`, orgID, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []MaintenanceTask
	for rows.Next() {
		var task MaintenanceTask
		if err := rows.Scan(&task.ID, &task.OrganizationID, &task.MaintenancePlanID, &task.AssetID, &task.ScheduledDate, &task.Status, &task.CompletedDate, &task.Notes, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) CreateMaintenanceTask(task *MaintenanceTask) error {
	result, err := r.Exec(`INSERT INTO maintenance_tasks (organization_id, maintenance_plan_id, asset_id, scheduled_date, status, notes) VALUES (?, ?, ?, ?, ?, ?)`,
		task.OrganizationID, task.MaintenancePlanID, task.AssetID, task.ScheduledDate, task.Status, task.Notes)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	task.ID = uint(id)
	return nil
}

func (r *Repository) UpdateMaintenanceTask(task *MaintenanceTask) error {
	_, err := r.Exec(`UPDATE maintenance_tasks SET status = ?, completed_date = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		task.Status, task.CompletedDate, task.Notes, task.ID)
	return err
}

func (r *Repository) GetMaintenancePlansDue(orgID uint) ([]MaintenancePlan, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := r.Query(`SELECT id, organization_id, asset_id, frequency_days, estimated_duration_hours, assigned_role, last_maintenance_date, next_maintenance_date, created_at, updated_at 
		FROM maintenance_plans WHERE organization_id = ? AND next_maintenance_date <= ?`, orgID, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []MaintenancePlan
	for rows.Next() {
		var plan MaintenancePlan
		if err := rows.Scan(&plan.ID, &plan.OrganizationID, &plan.AssetID, &plan.FrequencyDays, &plan.EstimatedDurationHours, &plan.AssignedRole, &plan.LastMaintenanceDate, &plan.NextMaintenanceDate, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *Repository) UpdateMaintenancePlanNextDate(planID uint, lastDate, nextDate time.Time) error {
	_, err := r.Exec(`UPDATE maintenance_plans SET last_maintenance_date = ?, next_maintenance_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, lastDate, nextDate, planID)
	return err
}

func (r *Repository) CreateWorkOrder(wo *WorkOrder) error {
	result, err := r.Exec(`INSERT INTO work_orders (organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, notes, created_by) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
	err := r.QueryRow(`SELECT id, organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, actual_start, actual_end, total_cost, notes, created_by, created_at, updated_at 
		FROM work_orders WHERE id = ? AND organization_id = ?`, id, orgID).
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
	err := r.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query = `SELECT id, organization_id, asset_id, technician_id, title, description, status, priority, scheduled_start, scheduled_end, actual_start, actual_end, total_cost, notes, created_by, created_at, updated_at 
		FROM work_orders WHERE organization_id = ?`
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

func (r *Repository) CreateWorkOrderPart(wop *WorkOrderPart) error {
	result, err := r.Exec(`INSERT INTO work_order_parts (work_order_id, part_id, quantity, unit_price, total_price) VALUES (?, ?, ?, ?, ?)`,
		wop.WorkOrderID, wop.PartID, wop.Quantity, wop.UnitPrice, wop.TotalPrice)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	wop.ID = uint(id)
	return nil
}

func (r *Repository) GetWorkOrderParts(woID uint) ([]WorkOrderPart, error) {
	rows, err := r.Query(`SELECT id, work_order_id, part_id, quantity, unit_price, total_price, created_at FROM work_order_parts WHERE work_order_id = ?`, woID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []WorkOrderPart
	for rows.Next() {
		var part WorkOrderPart
		if err := rows.Scan(&part.ID, &part.WorkOrderID, &part.PartID, &part.Quantity, &part.UnitPrice, &part.TotalPrice, &part.CreatedAt); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return parts, nil
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
	err := r.QueryRow(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, deleted_at, created_at, updated_at 
		FROM inventory_parts WHERE id = ? AND organization_id = ? AND deleted_at IS NULL`, id, orgID).
		Scan(&part.ID, &part.OrganizationID, &part.Name, &part.SKU, &part.Quantity, &part.MinThreshold, &part.CostPerUnit, &part.Location, &part.DeletedAt, &part.CreatedAt, &part.UpdatedAt)
	return part, err
}

func (r *Repository) ListInventoryParts(orgID uint, page, pageSize int) ([]InventoryPart, int, error) {
	offset := (page - 1) * pageSize

	var count int
	err := r.QueryRow(`SELECT COUNT(*) FROM inventory_parts WHERE organization_id = ? AND deleted_at IS NULL`, orgID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.Query(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, created_at, updated_at 
		FROM inventory_parts WHERE organization_id = ? AND deleted_at IS NULL ORDER BY name ASC LIMIT ? OFFSET ?`, orgID, pageSize, offset)
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

func (r *Repository) GetLowStockParts(orgID uint) ([]InventoryPart, error) {
	rows, err := r.Query(`SELECT id, organization_id, name, sku, quantity, min_threshold, cost_per_unit, location, created_at, updated_at 
		FROM inventory_parts WHERE organization_id = ? AND quantity <= min_threshold AND deleted_at IS NULL`, orgID)
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

func (r *Repository) UpdateInventoryPart(part *InventoryPart) error {
	_, err := r.Exec(`UPDATE inventory_parts SET name = ?, sku = ?, quantity = ?, min_threshold = ?, cost_per_unit = ?, location = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`,
		part.Name, part.SKU, part.Quantity, part.MinThreshold, part.CostPerUnit, part.Location, part.ID, part.OrganizationID)
	return err
}

func (r *Repository) DeductInventory(partID, orgID, quantity int) (int, error) {
	tx, err := r.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRow(`SELECT quantity FROM inventory_parts WHERE id = ? AND organization_id = ? AND deleted_at IS NULL`, partID, orgID).Scan(&currentQty)
	if err != nil {
		return 0, err
	}

	if currentQty < quantity {
		return 0, fmt.Errorf("insufficient stock")
	}

	_, err = tx.Exec(`UPDATE inventory_parts SET quantity = quantity - ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`, quantity, partID, orgID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return currentQty - quantity, nil
}

func (r *Repository) DeleteInventoryPart(id, orgID uint) error {
	_, err := r.Exec(`UPDATE inventory_parts SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`, id, orgID)
	return err
}

func (r *Repository) CreateDepreciation(depr *AssetDepreciation) error {
	result, err := r.Exec(`INSERT INTO asset_depreciation (organization_id, asset_id, year, depreciation_amount) VALUES (?, ?, ?, ?)`,
		depr.OrganizationID, depr.AssetID, depr.Year, depr.DepreciationAmount)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			_, err = r.Exec(`UPDATE asset_depreciation SET depreciation_amount = ? WHERE asset_id = ? AND year = ?`, depr.DepreciationAmount, depr.AssetID, depr.Year)
			return err
		}
		return err
	}
	id, _ := result.LastInsertId()
	depr.ID = uint(id)
	return nil
}

func (r *Repository) GetAssetDepreciation(assetID, orgID uint) ([]AssetDepreciation, error) {
	rows, err := r.Query(`SELECT id, organization_id, asset_id, year, depreciation_amount, created_at FROM asset_depreciation WHERE asset_id = ? AND organization_id = ? ORDER BY year ASC`, assetID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deprs []AssetDepreciation
	for rows.Next() {
		var depr AssetDepreciation
		if err := rows.Scan(&depr.ID, &depr.OrganizationID, &depr.AssetID, &depr.Year, &depr.DepreciationAmount, &depr.CreatedAt); err != nil {
			return nil, err
		}
		deprs = append(deprs, depr)
	}
	return deprs, nil
}

func (r *Repository) GetAssetMaintenanceCost(assetID, orgID uint) (float64, error) {
	var totalCost float64
	err := r.QueryRow(`SELECT COALESCE(SUM(total_cost), 0) FROM work_orders WHERE asset_id = ? AND organization_id = ?`, assetID, orgID).Scan(&totalCost)
	return totalCost, err
}

func (r *Repository) GetAllAssetCosts(orgID uint) ([]struct {
	AssetID   uint
	AssetName string
	TotalCost float64
}, error) {
	rows, err := r.Query(`SELECT a.id, a.name, COALESCE(SUM(wo.total_cost), 0) as total_cost 
		FROM assets a LEFT JOIN work_orders wo ON a.id = wo.asset_id 
		WHERE a.organization_id = ? AND a.deleted_at IS NULL 
		GROUP BY a.id, a.name`, orgID)
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

func (r *Repository) CreateAuditLog(log *AuditLog) error {
	result, err := r.Exec(`INSERT INTO audit_logs (organization_id, user_id, table_name, record_id, action, old_values, new_values) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.OrganizationID, log.UserID, log.TableName, log.RecordID, log.Action, log.OldValues, log.NewValues)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	log.ID = uint(id)
	return nil
}

func (r *Repository) GetAuditLogs(orgID uint, page, pageSize int, tableName string) ([]AuditLog, int, error) {
	offset := (page - 1) * pageSize

	var count int
	query := `SELECT COUNT(*) FROM audit_logs WHERE organization_id = ?`
	args := []interface{}{orgID}
	if tableName != "" {
		query += ` AND table_name = ?`
		args = append(args, tableName)
	}
	err := r.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query = `SELECT id, organization_id, user_id, table_name, record_id, action, old_values, new_values, created_at FROM audit_logs WHERE organization_id = ?`
	args = []interface{}{orgID}
	if tableName != "" {
		query += ` AND table_name = ?`
		args = append(args, tableName)
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(&log.ID, &log.OrganizationID, &log.UserID, &log.TableName, &log.RecordID, &log.Action, &log.OldValues, &log.NewValues, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, count, nil
}

func (r *Repository) GetDashboardStats(orgID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var assetCount int
	err := r.QueryRow(`SELECT COUNT(*) FROM assets WHERE organization_id = ? AND deleted_at IS NULL`, orgID).Scan(&assetCount)
	if err != nil {
		return nil, err
	}
	stats["asset_count"] = assetCount

	today := time.Now().Format("2006-01-02")
	var overdueCount int
	err = r.QueryRow(`SELECT COUNT(*) FROM maintenance_tasks mt 
		JOIN maintenance_plans mp ON mt.maintenance_plan_id = mp.id 
		WHERE mp.organization_id = ? AND mt.scheduled_date < ? AND mt.status IN ('pending', 'in_progress')`, orgID, today).Scan(&overdueCount)
	if err != nil {
		return nil, err
	}
	stats["overdue_maintenance"] = overdueCount

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

func (r *Repository) GetTechnicians(orgID uint) ([]User, error) {
	rows, err := r.Query(`SELECT id, organization_id, email, full_name, role, created_at, updated_at FROM users WHERE organization_id = ? AND role = 'technician'`, orgID)
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
