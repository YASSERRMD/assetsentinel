package worker

import (
	"assetsentinel/internal/repository"
	"assetsentinel/internal/websocket"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	repo    *repository.Repository
	hub     *websocket.Hub
	stop    chan bool
	running bool
	mu      sync.Mutex
}

func NewScheduler(repo *repository.Repository, hub *websocket.Hub) *Scheduler {
	return &Scheduler{
		repo: repo,
		hub:  hub,
		stop: make(chan bool),
	}
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkMaintenanceDue()
			s.checkOverdueTasks()
		case <-s.stop:
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.stop <- true
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

func (s *Scheduler) checkMaintenanceDue() {
	orgs, err := s.repo.ListOrganizations()
	if err != nil {
		log.Printf("Error fetching organizations: %v", err)
		return
	}

	for _, org := range orgs {
		plans, err := s.repo.GetMaintenancePlansDue(org.ID)
		if err != nil {
			log.Printf("Error fetching maintenance plans for org %d: %v", org.ID, err)
			continue
		}

		for _, plan := range plans {
			task := &repository.MaintenanceTask{
				OrganizationID:    plan.OrganizationID,
				MaintenancePlanID: plan.ID,
				AssetID:           plan.AssetID,
				ScheduledDate:     plan.NextMaintenanceDate,
				Status:            "pending",
			}

			if err := s.repo.CreateMaintenanceTask(task); err != nil {
				log.Printf("Error creating maintenance task: %v", err)
				continue
			}

			s.hub.BroadcastToOrg(org.ID, map[string]interface{}{
				"type":           "maintenance_due",
				"maintenance_id": plan.ID,
				"asset_id":       plan.AssetID,
				"scheduled_date": plan.NextMaintenanceDate,
			})
		}
	}
}

func (s *Scheduler) checkOverdueTasks() {
	orgs, err := s.repo.ListOrganizations()
	if err != nil {
		log.Printf("Error fetching organizations: %v", err)
		return
	}

	for _, org := range orgs {
		tasks, err := s.repo.GetOverdueMaintenanceTasks(org.ID)
		if err != nil {
			log.Printf("Error fetching overdue tasks for org %d: %v", org.ID, err)
			continue
		}

		for _, task := range tasks {
			task.Status = "overdue"
			if err := s.repo.UpdateMaintenanceTask(&task); err != nil {
				log.Printf("Error updating overdue task: %v", err)
				continue
			}

			s.hub.BroadcastToOrg(org.ID, map[string]interface{}{
				"type":           "maintenance_overdue",
				"maintenance_id": task.ID,
				"asset_id":       task.AssetID,
				"scheduled_date": task.ScheduledDate,
			})
		}
	}
}
