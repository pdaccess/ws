package service

import (
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/database"
)

type Impl struct {
	inventoryRepo *database.InventoryRepository
}

func New(inventoryRepo *database.InventoryRepository) ports.Service {
	return &Impl{
		inventoryRepo: inventoryRepo,
	}
}
