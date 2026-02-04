package service

import (
	"git.h2hsecure.com/core/ws/internal/core/ports"
)

type Impl struct {
}

func New() ports.Service {
	return &Impl{}
}
