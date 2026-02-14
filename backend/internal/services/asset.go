package services

import (
	"assetsentinel/internal/repository"
)

type AssetService struct {
	repo *repository.Repository
}

func NewAssetService(repo *repository.Repository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) Create(asset *repository.Asset) error {
	return s.repo.CreateAsset(asset)
}

func (s *AssetService) Get(id, orgID uint) (*repository.Asset, error) {
	return s.repo.GetAsset(id, orgID)
}

func (s *AssetService) List(orgID uint, page, pageSize int, status, category string) ([]repository.Asset, int, error) {
	return s.repo.ListAssets(orgID, page, pageSize, status, category)
}

func (s *AssetService) Update(asset *repository.Asset) error {
	return s.repo.UpdateAsset(asset)
}

func (s *AssetService) Delete(id, orgID uint) error {
	return s.repo.SoftDeleteAsset(id, orgID)
}
