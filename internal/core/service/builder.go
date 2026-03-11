package service

import (
	"git.h2hsecure.com/core/ws/internal/core/ports"
	"git.h2hsecure.com/core/ws/internal/database"
)

type Impl struct {
	inventoryRepo *database.InventoryRepository
}

func New(inventoryRepo *database.InventoryRepository) ports.Service {
	return &Impl{
		inventoryRepo: inventoryRepo,
	}
}
