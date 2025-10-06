package snrdev

import (
	"context"
	"strconv"
	"strings"

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

func (u *Onboarding) getAppmans(memberId string, ctx context.Context, tx *gorm.DB) (appman []AppmanDB, err error) {
	var scopeAppman []func(*gorm.DB) *gorm.DB

	scopeAppman = append(scopeAppman, Where("member_id = ?", memberId))
	scopeAppman = append(scopeAppman, Where("dopa_status = ?", true))
	scopeAppman = append(scopeAppman, Order("updated_at desc"))
	err = u.repo.Appman.GetAll(ctx, u.db, &appman, scopeAppman...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) buildContactScope(mobile string, email string) (scope []func(*gorm.DB) *gorm.DB) {
	mobile = strings.TrimSpace(mobile)
	email = strings.TrimSpace(email)
	switch {
	case email != "" && mobile != "":
		scope = append(scope, Where("email = ? OR mobile = ?", email, mobile))
	case email != "":
		scope = append(scope, Where("email = ?", email))
	case mobile != "":
		scope = append(scope, Where("mobile = ?", mobile))
	}
	return
}

func (u *Onboarding) upsertPreMember(emailData string, mobileData string, memberId string, appmans []AppmanDB, ctx context.Context, tx *gorm.DB) (err error) {
	var scope []func(*gorm.DB) *gorm.DB

	if len(appmans) == 0 {
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
		scope = append(scope, u.buildContactScope(emailData, mobileData)...)
		if err = u.repo.PreMember.GetAll(ctx, u.db, &queryPremember, scope...); utils.RollbackOnError(tx, err) {
			return
		}
		if len(queryPremember) <= 0 {
			if err = u.repo.PreMember.CreateAll(ctx, u.db, &premember); utils.RollbackOnError(tx, err) {
				return
			}
		} else {
			if err = u.repo.PreMember.UpdateAll(ctx, tx, &premember[0], scope...); utils.RollbackOnError(tx, err) {
				return
			}
		}
	}
	return
}

func (u *Onboarding) getJoinedCustomer(memberId string, ctx context.Context, tx *gorm.DB) (preMember []PreMemberDB, err error) {
	var scopeCustomerData []func(*gorm.DB) *gorm.DB

	scopeCustomerData = append(scopeCustomerData, Where("member_id = ?", memberId))
	err = u.repo.PreMember.JoinCustomer(ctx, u.db, &preMember, scopeCustomerData...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) mapToMemberResponse(preMember []PreMemberDB, emailData string, mobileData string, memberId string) (step int, resp PreCitizenshipResp) {
	if len(preMember) > 0 {
		step = preMember[0].Customer.Step
		resp.CustomerData = CustomerData{
			Email:    preMember[0].Email,
			Mobile:   preMember[0].Mobile,
			MemberID: preMember[0].MemberID,
			Fullname: &CustomerFullname{
				ID:           preMember[0].Customer.ID,
				MemberID:     preMember[0].Customer.MemberID,
				Citizenship:  preMember[0].Customer.Citizenship,
				Title:        preMember[0].Customer.Title,
				ThName:       preMember[0].Customer.ThName,
				ThMiddleName: preMember[0].Customer.ThMiddleName,
				ThSurname:    preMember[0].Customer.ThSurname,
				EnName:       preMember[0].Customer.EnName,
				EnMiddleName: preMember[0].Customer.EnMiddleName,
				EnSurname:    preMember[0].Customer.EnSurname,
				Mobile:       preMember[0].Customer.Mobile,
				Email:        preMember[0].Customer.Email,
				Agreement:    preMember[0].Customer.Agreement,
				Step:         preMember[0].Customer.Step,
				Type:         preMember[0].Customer.Type,
			},
			IDCard: &IDCardDetail{
				DateOfBirth: preMember[0].Customer.DateOfBirth,
				Status:      preMember[0].Customer.Status,
				IDCard:      preMember[0].Customer.IDCard,
				LaserCode:   preMember[0].Customer.LaserCode,
				ExpireDate:  preMember[0].Customer.ExpireDate,
			},
			SuiteTest:     &preMember[0].SuiteTest,
			Document:      &preMember[0].Document,
			Addresses:     preMember[0].Addresses,
			SourceOfFund:  &preMember[0].SourceOfFund,
			Occupation:    &preMember[0].Occupation,
			Banks:         preMember[0].Banks,
			KnowledgeTest: preMember[0].Customer.KnowledgeTest,
		}

		if preMember[0].SuiteTest.ID == "" {
			resp.CustomerData.SuiteTest = nil
		}
		if preMember[0].Document.ID == "" {
			resp.CustomerData.Document = nil
		}
		if preMember[0].SourceOfFund.ID == "" {
			resp.CustomerData.SourceOfFund = nil
		}
		if preMember[0].Occupation.ID == "" {
			resp.CustomerData.Occupation = nil
		}
		if len(preMember[0].Banks) == 0 {
			resp.CustomerData.Banks = nil
		}
		if len(preMember[0].Addresses) == 0 {
			resp.CustomerData.Addresses = nil
		}
	}

	resp.Member.ID = memberId
	resp.Member.Email = emailData
	resp.Member.Mobile = mobileData

	return
}

func (u *Onboarding) getCustomer(memberId string, ctx context.Context, tx *gorm.DB) (customer []CustomerDetailsDB, err error) {
	var scopeCustomer []func(*gorm.DB) *gorm.DB
	scopeCustomer = append(scopeCustomer, Where("member_id = ?", memberId))
	err = u.repo.Customer.GetAll(ctx, u.db, &customer, scopeCustomer...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) setCitizenship(resp *PreCitizenshipResp, memberId string, ctx context.Context, tx *gorm.DB) (err error) {
	resp.Member.Citizenship = 1
	customers, err := u.getCustomer(memberId, ctx, tx)
	if err != nil {
		return
	}
	if len(customers) > 0 {
		resp.Member.Citizenship = customers[0].Citizenship
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
	appmans, err := u.getAppmans(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	if err = u.upsertPreMember(data.Email, data.Mobile, memberId, appmans, ctx, tx); err != nil {
		return PreCitizenshipResp{}, err
	}

	preMember, err := u.getJoinedCustomer(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	step, resp := u.mapToMemberResponse(preMember, data.Email, data.Mobile, memberId)
	if len(appmans) > 0 {
		if step == 0 {
			step = 50
		}
		resp.Step = step
		u.setCitizenship(&resp, memberId, ctx, tx)
	}

	return
}
