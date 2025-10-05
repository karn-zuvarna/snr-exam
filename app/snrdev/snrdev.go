package snrdev

import (
	"context"
	"log"
	"strconv"

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

func getMemberId(memberID int, userID int) string {
	if memberID > 0 {
		return strconv.Itoa(memberID)
	} else {
		return strconv.Itoa(userID)
	}
}

func getScopeAppman(scopeAppman []func(*gorm.DB) *gorm.DB, memberId string, appmanRepo IAppman, ctx context.Context, db *gorm.DB, tx *gorm.DB) ([]AppmanDB, error) {
	var appman []AppmanDB

	scopeAppman = append(scopeAppman, Where("member_id = ?", memberId))
	scopeAppman = append(scopeAppman, Where("dopa_status = ?", true))
	scopeAppman = append(scopeAppman, Order("updated_at desc"))
	err := appmanRepo.GetAll(ctx, db, &appman, scopeAppman...)
	if err != nil {
		tx.Rollback()
	}
	return appman, err
}

func (u *Onboarding) Snr_Dev_Exam(ctx context.Context, token string, publicKey string) (resp PreCitizenshipResp, err error) {
	var scope, scopeAppman, scopeCustomer, scopeCustomerData []func(*gorm.DB) *gorm.DB
	var step int

	tx := u.db.Begin()
	defer tx.Commit()

	data, err := u.service.VerifyToken(token, []byte(publicKey))
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	memberId := getMemberId(data.MemberID, data.UserID)
	var appman []AppmanDB
	appman, err = getScopeAppman(scopeAppman, memberId, u.repo.Appman, ctx, u.db, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	if len(appman) == 0 {
		var queryPremember = []PreMemberDB{}
		premember := []PreMemberDB{
			{
				ID:       u.ext.GenUuid(),
				Email:    data.Email,
				Mobile:   data.Mobile,
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
		if len(data.Mobile) > 0 && len(data.Email) > 0 {
			scope = append(scope, Where("email = ? OR mobile = ?", data.Email, data.Mobile))
		}
		if len(data.Mobile) == 0 && len(data.Email) > 0 {
			scope = append(scope, Where("email = ?", data.Email))
		}
		if len(data.Mobile) > 0 && len(data.Email) == 0 {
			scope = append(scope, Where("mobile = ?", data.Mobile))
		}
		if err = u.repo.PreMember.GetAll(ctx, u.db, &queryPremember, scope...); err != nil {
			tx.Rollback()
			return PreCitizenshipResp{}, err
		}
		if len(queryPremember) <= 0 {
			if err = u.repo.PreMember.CreateAll(ctx, u.db, &premember); err != nil {
				tx.Rollback()
				return PreCitizenshipResp{}, err
			}
		} else {
			if err = u.repo.PreMember.UpdateAll(ctx, tx, &premember[0], scope...); err != nil {
				tx.Rollback()
				return PreCitizenshipResp{}, err
			}
		}
	}
	var preMemberResp []PreMemberDB
	scopeCustomerData = append(scopeCustomerData, Where("member_id = ?", memberId))

	err = u.repo.PreMember.JoinCustomer(ctx, u.db, &preMemberResp, scopeCustomerData...)
	if err != nil {
		tx.Rollback()
		return PreCitizenshipResp{}, err
	}
	log.Printf("preMemberResp: %+v", preMemberResp)
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
	resp.Member.Email = data.Email
	resp.Member.Mobile = data.Mobile
	if len(appman) > 0 {
		if step == 0 {
			step = 50
		}
		resp.Member.Citizenship = 1
		var customer []CustomerDetailsDB
		scopeCustomer = append(scopeCustomer, Where("member_id = ?", memberId))
		if err = u.repo.Customer.GetAll(ctx, u.db, &customer, scopeCustomer...); err != nil {
			tx.Rollback()
			return PreCitizenshipResp{}, err
		}
		if len(customer) > 0 {
			resp.Member.Citizenship = customer[0].Citizenship
		}
		resp.Step = step
	}

	return
}
