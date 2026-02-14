package repository

import (
	"time"
)

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
