package service

import (
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/database"
)

type Impl struct {
	assetRepo *database.AssetRepository
	auditRepo *database.AuditRepository
	vecGen    ports.VectorGenerator
}

func New(assetRepo *database.AssetRepository, auditRepo *database.AuditRepository, vecGen ports.VectorGenerator) ports.Service {
	return &Impl{
		assetRepo: assetRepo,
		auditRepo: auditRepo,
		vecGen:    vecGen,
	}
}
