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
