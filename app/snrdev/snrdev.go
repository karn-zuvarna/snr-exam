package snrdev

import (
	"context"
	"strconv"

	"github.com/karn-zuvarna/snr-exam/app/utils"

	"gorm.io/gorm"
)

type Onboarding struct {
	ext     IExternal
	db      *gorm.DB
	service IService
	repo    Irepo
}

func New(ext IExternal, db *gorm.DB, service IService, repo Irepo) *Onboarding {
	return &Onboarding{
		ext:     ext,
		db:      db,
		service: service,
		repo:    repo,
	}
}

func (u *Onboarding) getMemberId(memberID int, userID int) string {
	if memberID > 0 {
		return strconv.Itoa(memberID)
	} else {
		return strconv.Itoa(userID)
	}
}

func (u *Onboarding) getScopeAppman(memberId string, ctx context.Context, tx *gorm.DB) (appman []AppmanDB, err error) {
	var scopeAppman []func(*gorm.DB) *gorm.DB

	scopeAppman = append(scopeAppman, Where("member_id = ?", memberId))
	scopeAppman = append(scopeAppman, Where("dopa_status = ?", true))
	scopeAppman = append(scopeAppman, Order("updated_at desc"))
	err = u.repo.Appman.GetAll(ctx, u.db, &appman, scopeAppman...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) upsertPreMember(emailData string, mobileData string, memberId string, appman []AppmanDB, ctx context.Context, tx *gorm.DB) (err error) {
	var scope []func(*gorm.DB) *gorm.DB

	if len(appman) == 0 {
		var queryPremember = []PreMemberDB{}
		premember := []PreMemberDB{
			{
				ID:       u.ext.GenUuid(),
				Email:    emailData,
				Mobile:   mobileData,
				MemberID: memberId,
			},
		}

		col := []string{
			"created_at",
			"updated_at",
			"id",
			"member_id",
			"email",
			"mobile",
		}

		scope = append(scope, Select(col))
		if len(mobileData) > 0 && len(emailData) > 0 {
			scope = append(scope, Where("email = ? OR mobile = ?", emailData, mobileData))
		}
		if len(mobileData) == 0 && len(emailData) > 0 {
			scope = append(scope, Where("email = ?", emailData))
		}
		if len(mobileData) > 0 && len(emailData) == 0 {
			scope = append(scope, Where("mobile = ?", mobileData))
		}
		if err := u.repo.PreMember.GetAll(ctx, u.db, &queryPremember, scope...); utils.RollbackOnError(tx, err) {
			return err
		}
		if len(queryPremember) <= 0 {
			if err := u.repo.PreMember.CreateAll(ctx, u.db, &premember); utils.RollbackOnError(tx, err) {
				return err
			}
		} else {
			if err := u.repo.PreMember.UpdateAll(ctx, tx, &premember[0], scope...); utils.RollbackOnError(tx, err) {
				return err
			}
		}
	}
	return
}

func (u *Onboarding) getJoinedCustomerInfo(memberId string, ctx context.Context, tx *gorm.DB) (preMemberResp []PreMemberDB, err error) {
	var scopeCustomerData []func(*gorm.DB) *gorm.DB

	scopeCustomerData = append(scopeCustomerData, Where("member_id = ?", memberId))
	err = u.repo.PreMember.JoinCustomer(ctx, u.db, &preMemberResp, scopeCustomerData...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) mapToMemberResponse(preMemberResp []PreMemberDB, emailData string, mobileData string, memberId string) (step int, resp PreCitizenshipResp) {
	if len(preMemberResp) > 0 {
		step = preMemberResp[0].Customer.Step
		resp.CustomerData = CustomerData{
			Email:    preMemberResp[0].Email,
			Mobile:   preMemberResp[0].Mobile,
			MemberID: preMemberResp[0].MemberID,
			Fullname: &CustomerFullname{
				ID:           preMemberResp[0].Customer.ID,
				MemberID:     preMemberResp[0].Customer.MemberID,
				Citizenship:  preMemberResp[0].Customer.Citizenship,
				Title:        preMemberResp[0].Customer.Title,
				ThName:       preMemberResp[0].Customer.ThName,
				ThMiddleName: preMemberResp[0].Customer.ThMiddleName,
				ThSurname:    preMemberResp[0].Customer.ThSurname,
				EnName:       preMemberResp[0].Customer.EnName,
				EnMiddleName: preMemberResp[0].Customer.EnMiddleName,
				EnSurname:    preMemberResp[0].Customer.EnSurname,
				Mobile:       preMemberResp[0].Customer.Mobile,
				Email:        preMemberResp[0].Customer.Email,
				Agreement:    preMemberResp[0].Customer.Agreement,
				Step:         preMemberResp[0].Customer.Step,
				Type:         preMemberResp[0].Customer.Type,
			},
			IDCard: &IDCardDetail{
				DateOfBirth: preMemberResp[0].Customer.DateOfBirth,
				Status:      preMemberResp[0].Customer.Status,
				IDCard:      preMemberResp[0].Customer.IDCard,
				LaserCode:   preMemberResp[0].Customer.LaserCode,
				ExpireDate:  preMemberResp[0].Customer.ExpireDate,
			},
			SuiteTest:     &preMemberResp[0].SuiteTest,
			Document:      &preMemberResp[0].Document,
			Addresses:     preMemberResp[0].Addresses,
			SourceOfFund:  &preMemberResp[0].SourceOfFund,
			Occupation:    &preMemberResp[0].Occupation,
			Banks:         preMemberResp[0].Banks,
			KnowledgeTest: preMemberResp[0].Customer.KnowledgeTest,
		}

		if preMemberResp[0].SuiteTest.ID == "" {
			resp.CustomerData.SuiteTest = nil
		}
		if preMemberResp[0].Document.ID == "" {
			resp.CustomerData.Document = nil
		}
		if preMemberResp[0].SourceOfFund.ID == "" {
			resp.CustomerData.SourceOfFund = nil
		}
		if preMemberResp[0].Occupation.ID == "" {
			resp.CustomerData.Occupation = nil
		}
		if len(preMemberResp[0].Banks) == 0 {
			resp.CustomerData.Banks = nil
		}
		if len(preMemberResp[0].Addresses) == 0 {
			resp.CustomerData.Addresses = nil
		}
	}

	resp.Member.ID = memberId
	resp.Member.Email = emailData
	resp.Member.Mobile = mobileData

	return
}

func (u *Onboarding) getCustomerInfo(memberId string, ctx context.Context, tx *gorm.DB) (customer []CustomerDetailsDB, err error) {
	var scopeCustomer []func(*gorm.DB) *gorm.DB
	scopeCustomer = append(scopeCustomer, Where("member_id = ?", memberId))
	err = u.repo.Customer.GetAll(ctx, u.db, &customer, scopeCustomer...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) setCitizenship(resp *PreCitizenshipResp, memberId string, ctx context.Context, tx *gorm.DB) (err error) {
	resp.Member.Citizenship = 1
	customer, err := u.getCustomerInfo(memberId, ctx, tx)
	if err != nil {
		return err
	}
	if len(customer) > 0 {
		resp.Member.Citizenship = customer[0].Citizenship
	}
	return
}

func (u *Onboarding) Snr_Dev_Exam_v1(ctx context.Context, token string, publicKey string) (resp PreCitizenshipResp, err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	data, err := u.service.VerifyToken(token, []byte(publicKey))
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	memberId := u.getMemberId(data.MemberID, data.UserID)
	appman, err := u.getScopeAppman(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	if err = u.upsertPreMember(data.Email, data.Mobile, memberId, appman, ctx, tx); err != nil {
		return PreCitizenshipResp{}, err
	}

	preMemberResp, err := u.getJoinedCustomerInfo(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	step, resp := u.mapToMemberResponse(preMemberResp, data.Email, data.Mobile, memberId)
	if len(appman) > 0 {
		if step == 0 {
			step = 50
		}
		resp.Step = step
		u.setCitizenship(&resp, memberId, ctx, tx)
	}

	return
}
