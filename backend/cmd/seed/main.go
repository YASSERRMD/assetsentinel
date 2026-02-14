package main

import (
	"assetsentinel/internal/config"
	"assetsentinel/internal/repository"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}

	repo := repository.NewRepository(db)

	org := &repository.Organization{Name: "Acme Corp"}
	if err := repo.CreateOrganization(org); err != nil {
		log.Printf("Org exists: %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	admin := &repository.User{
		OrganizationID: org.ID,
		Email:          "admin@acme.com",
		PasswordHash:   string(hash),
		FullName:       "Admin User",
		Role:           "admin",
	}
	repo.CreateUser(admin)

	log.Println("Seed data created!")
}
