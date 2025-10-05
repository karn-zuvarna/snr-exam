package snrdev

import (
	"context"

	"gorm.io/gorm"
)

func (u *Onboarding) TestGetScopeAppman(scopeAppman []func(*gorm.DB) *gorm.DB, memberId string, ctx context.Context, tx *gorm.DB) ([]AppmanDB, error) {
	return u.getScopeAppman(scopeAppman, memberId, ctx, tx)
}
