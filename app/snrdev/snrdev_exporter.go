package snrdev

import (
	"context"

	"gorm.io/gorm"
)

func (u *Onboarding) TestGetScopeAppman(scopeAppman []func(*gorm.DB) *gorm.DB, memberId string, ctx context.Context, tx *gorm.DB) ([]AppmanDB, error) {
	return u.getScopeAppman(scopeAppman, memberId, ctx, tx)
}

func (u *Onboarding) TestUpsertPremember(uuid string, emailData string, mobileData string, memberId string, appman []AppmanDB, ctx context.Context, tx *gorm.DB) error {
	return u.upsertPreMember(uuid, emailData, mobileData, memberId, appman, ctx, tx)
}
