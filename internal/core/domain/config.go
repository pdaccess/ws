package domain

import (
	"fmt"

	"github.com/samber/lo"
)

const (
	SmtpConfigContext   = ItemContext("smtp")
	PortalConfigContext = ItemContext("portal")
)

var (
	ValidConfigContexts = []ItemContext{
		SmtpConfigContext, PortalConfigContext,
	}
)

type ItemContext string

type ConfigItem struct {
	Key, Value string
}

func (i ItemContext) Validate() error {
	if !lo.Contains(ValidConfigContexts, i) {
		return fmt.Errorf("wrong config context: %s", i)
	}

	return nil
}
