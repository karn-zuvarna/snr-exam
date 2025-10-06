package snrdev

import (
	"context"

	"gorm.io/gorm"
)

func (u *Onboarding) TestGetMemberId(memberID int, userID int) string {
	return u.getMemberId(memberID, userID)
}

func (u *Onboarding) TestGetAppman(memberId string, ctx context.Context, tx *gorm.DB) ([]AppmanDB, error) {
	return u.getAppman(memberId, ctx, tx)
}

func (u *Onboarding) TestBuildContactScope(mobile string, email string) (scope []func(*gorm.DB) *gorm.DB) {
	return u.buildContactScope(mobile, email)
}

func (u *Onboarding) TestUpsertPremember(emailData string, mobileData string, memberId string, appman []AppmanDB, ctx context.Context, tx *gorm.DB) error {
	return u.upsertPreMember(emailData, mobileData, memberId, appman, ctx, tx)
}

func (u *Onboarding) TestGetJoinedCustomer(memberId string, ctx context.Context, tx *gorm.DB) (preMemberResp []PreMemberDB, err error) {
	return u.getJoinedCustomer(memberId, ctx, tx)
}

func (u *Onboarding) TestMapToPreMemberResponse(preMemberResp []PreMemberDB, emailData string, mobileData string, memberId string) (step int, resp PreCitizenshipResp) {
	return u.mapToPreMemberResponse(preMemberResp, emailData, mobileData, memberId)
}

func (u *Onboarding) TestGetCustomer(memberId string, ctx context.Context, tx *gorm.DB) (customer []CustomerDetailsDB, err error) {
	return u.getCustomer(memberId, ctx, tx)
}

func (u *Onboarding) TestSetCitizenship(resp *PreCitizenshipResp, memberId string, ctx context.Context, tx *gorm.DB) (err error) {
	return u.setCitizenship(resp, memberId, ctx, tx)
}
