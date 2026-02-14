package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")

	return &DB{db}, nil
}

func RunMigrations(db *DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('admin', 'maintenance_manager', 'technician', 'viewer')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_org ON users(organization_id)`,

		`CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			serial_number TEXT,
			installation_date DATE,
			location TEXT,
			purchase_cost REAL DEFAULT 0,
			warranty_expiry DATE,
			status TEXT DEFAULT 'active' CHECK(status IN ('active', 'retired', 'under_maintenance')),
			deleted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_assets_org ON assets(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status)`,

		`CREATE TABLE IF NOT EXISTS maintenance_plans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			asset_id INTEGER NOT NULL,
			frequency_days INTEGER NOT NULL,
			estimated_duration_hours REAL,
			assigned_role TEXT CHECK(assigned_role IN ('technician', 'maintenance_manager')),
			last_maintenance_date DATE,
			next_maintenance_date DATE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mp_org ON maintenance_plans(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mp_next_date ON maintenance_plans(next_maintenance_date)`,

		`CREATE TABLE IF NOT EXISTS maintenance_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			maintenance_plan_id INTEGER NOT NULL,
			asset_id INTEGER NOT NULL,
			scheduled_date DATE NOT NULL,
			status TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'in_progress', 'completed', 'overdue')),
			completed_date DATE,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (maintenance_plan_id) REFERENCES maintenance_plans(id) ON DELETE CASCADE,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mt_org ON maintenance_tasks(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mt_status ON maintenance_tasks(status)`,

		`CREATE TABLE IF NOT EXISTS work_orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			asset_id INTEGER NOT NULL,
			technician_id INTEGER,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'in_progress', 'completed', 'closed')),
			priority TEXT DEFAULT 'medium' CHECK(priority IN ('low', 'medium', 'high', 'critical')),
			scheduled_start DATETIME,
			scheduled_end DATETIME,
			actual_start DATETIME,
			actual_end DATETIME,
			total_cost REAL DEFAULT 0,
			notes TEXT,
			created_by INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_wo_org ON work_orders(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_wo_status ON work_orders(status)`,

		`CREATE TABLE IF NOT EXISTS inventory_parts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			sku TEXT UNIQUE NOT NULL,
			quantity INTEGER DEFAULT 0,
			min_threshold INTEGER DEFAULT 0,
			cost_per_unit REAL DEFAULT 0,
			location TEXT,
			deleted_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_inv_org ON inventory_parts(organization_id)`,

		`CREATE TABLE IF NOT EXISTS asset_depreciation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			organization_id INTEGER NOT NULL,
			asset_id INTEGER NOT NULL,
			year INTEGER NOT NULL,
			depreciation_amount REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			UNIQUE(asset_id, year)
		)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
