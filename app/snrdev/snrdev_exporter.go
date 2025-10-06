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

func (u *Onboarding) TestUpsertPreMember(email string, mobile string, memberId string, appman []AppmanDB, ctx context.Context, tx *gorm.DB) error {
	return u.upsertPreMember(email, mobile, memberId, appman, ctx, tx)
}

func (u *Onboarding) TestGetPreMemberJoinedCustomer(memberId string, ctx context.Context, tx *gorm.DB) (preMember []PreMemberDB, err error) {
	return u.getPreMemberJoinedCustomer(memberId, ctx, tx)
}

func (u *Onboarding) TestConvertPreMemberToPreCitizenshipResp(preMembers []PreMemberDB, email string, mobile string, memberId string) (step int, resp PreCitizenshipResp) {
	return u.convertPreMemberToPreCitizenshipResp(preMembers, email, mobile, memberId)
}

func (u *Onboarding) TestGetCustomer(memberId string, ctx context.Context, tx *gorm.DB) (customer []CustomerDetailsDB, err error) {
	return u.getCustomer(memberId, ctx, tx)
}

func (u *Onboarding) TestSetCitizenship(resp *PreCitizenshipResp, memberId string, ctx context.Context, tx *gorm.DB) (err error) {
	return u.setCitizenship(resp, memberId, ctx, tx)
}
