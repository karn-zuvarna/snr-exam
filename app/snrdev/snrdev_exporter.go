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

func (u *Onboarding) TestGetJoinedCustomerInfo(memberId string, ctx context.Context, tx *gorm.DB) (preMemberResp []PreMemberDB, err error) {
	return u.getJoinedCustomerInfo(memberId, ctx, tx)
}

func (u *Onboarding) TestMapToMemberResponse(preMemberResp []PreMemberDB, emailData string, mobileData string, memberId string) (step int, resp PreCitizenshipResp) {
	return u.mapToMemberResponse(preMemberResp, emailData, mobileData, memberId)
}

func (u *Onboarding) TestGetCustomerInfo(memberId string, ctx context.Context, tx *gorm.DB) (customer []CustomerDetailsDB, err error) {
	return u.getCustomerInfo(memberId, ctx, tx)
}
