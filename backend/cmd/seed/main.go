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
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewRepository(db)

	org := &repository.Organization{Name: "Acme Corporation"}
	if err := repo.CreateOrganization(org); err != nil {
		log.Printf("Organization already exists or error: %v", err)
	} else {
		log.Printf("Created organization: %s (ID: %d)", org.Name, org.ID)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	admin := &repository.User{
		OrganizationID: org.ID,
		Email:          "admin@acme.com",
		PasswordHash:   string(hash),
		FullName:       "System Admin",
		Role:           "admin",
	}
	if err := repo.CreateUser(admin); err != nil {
		log.Printf("Admin user already exists or error: %v", err)
	} else {
		log.Printf("Created admin user: %s", admin.Email)
	}

	techHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	tech := &repository.User{
		OrganizationID: org.ID,
		Email:          "tech@acme.com",
		PasswordHash:   string(techHash),
		FullName:       "John Technician",
		Role:           "technician",
	}
	if err := repo.CreateUser(tech); err != nil {
		log.Printf("Tech user already exists or error: %v", err)
	} else {
		log.Printf("Created technician user: %s", tech.Email)
	}

	hvac001 := "HVAC-001"
	buildingA := "Building A"
	gen001 := "GEN-001"
	basement := "Basement"
	elv001 := "ELV-001"
	mainLobby := "Main Lobby"
	fp001 := "FP-001"
	ch001 := "CH-001"
	roof := "Roof"

	assets := []repository.Asset{
		{OrganizationID: org.ID, Name: "HVAC Unit 1", Category: "HVAC", SerialNumber: &hvac001, Location: &buildingA, PurchaseCost: 15000, Status: "active"},
		{OrganizationID: org.ID, Name: "Generator 1", Category: "Electrical", SerialNumber: &gen001, Location: &basement, PurchaseCost: 25000, Status: "active"},
		{OrganizationID: org.ID, Name: "Elevator 1", Category: "Transportation", SerialNumber: &elv001, Location: &mainLobby, PurchaseCost: 75000, Status: "active"},
		{OrganizationID: org.ID, Name: "Fire Pump 1", Category: "Safety", SerialNumber: &fp001, Location: &basement, PurchaseCost: 12000, Status: "active"},
		{OrganizationID: org.ID, Name: "Chiller 1", Category: "HVAC", SerialNumber: &ch001, Location: &roof, PurchaseCost: 45000, Status: "under_maintenance"},
	}
	for _, asset := range assets {
		if err := repo.CreateAsset(&asset); err != nil {
			log.Printf("Error creating asset %s: %v", asset.Name, err)
		} else {
			log.Printf("Created asset: %s", asset.Name)
		}
	}

	storageA := "Storage Room A"
	storageB := "Storage Room B"
	storageC := "Storage Room C"

	parts := []repository.InventoryPart{
		{OrganizationID: org.ID, Name: "Air Filter", SKU: "AF-001", Quantity: 50, MinThreshold: 10, CostPerUnit: 25.00, Location: &storageA},
		{OrganizationID: org.ID, Name: "Belt Drive", SKU: "BD-001", Quantity: 15, MinThreshold: 5, CostPerUnit: 45.00, Location: &storageA},
		{OrganizationID: org.ID, Name: "Contactor", SKU: "CT-001", Quantity: 8, MinThreshold: 10, CostPerUnit: 35.00, Location: &storageB},
		{OrganizationID: org.ID, Name: "Capacitor", SKU: "CP-001", Quantity: 25, MinThreshold: 15, CostPerUnit: 15.00, Location: &storageB},
		{OrganizationID: org.ID, Name: "Motor Oil", SKU: "MO-001", Quantity: 100, MinThreshold: 20, CostPerUnit: 12.00, Location: &storageC},
	}
	for _, part := range parts {
		if err := repo.CreateInventoryPart(&part); err != nil {
			log.Printf("Error creating part %s: %v", part.Name, err)
		} else {
			log.Printf("Created inventory part: %s", part.Name)
		}
	}

	log.Println("Seed data completed successfully!")
}
