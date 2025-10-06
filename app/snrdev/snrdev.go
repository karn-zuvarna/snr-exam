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

func (u *Onboarding) getAppman(memberId string, ctx context.Context, tx *gorm.DB) (appman []AppmanDB, err error) {
	var scope []func(*gorm.DB) *gorm.DB

	scope = append(scope, Where("member_id = ?", memberId))
	scope = append(scope, Where("dopa_status = ?", true))
	scope = append(scope, Order("updated_at desc"))
	err = u.repo.Appman.GetAll(ctx, u.db, &appman, scope...)
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

func (u *Onboarding) upsertPreMember(email string, mobile string, memberId string, appmans []AppmanDB, ctx context.Context, tx *gorm.DB) (err error) {
	var scope []func(*gorm.DB) *gorm.DB

	if len(appmans) == 0 {
		var queryPremember = []PreMemberDB{}
		premember := []PreMemberDB{
			{
				ID:       u.ext.GenUuid(),
				Email:    email,
				Mobile:   mobile,
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
		scope = append(scope, u.buildContactScope(email, mobile)...)
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

func (u *Onboarding) getPreMemberJoinedCustomer(memberId string, ctx context.Context, tx *gorm.DB) (preMember []PreMemberDB, err error) {
	var scopeCustomerData []func(*gorm.DB) *gorm.DB

	scopeCustomerData = append(scopeCustomerData, Where("member_id = ?", memberId))
	err = u.repo.PreMember.JoinCustomer(ctx, u.db, &preMember, scopeCustomerData...)
	utils.RollbackOnError(tx, err)
	return
}

func (u *Onboarding) mapToPreMemberResponse(preMembers []PreMemberDB, email string, mobile string, memberId string) (step int, resp PreCitizenshipResp) {
	if len(preMembers) > 0 {
		step = preMembers[0].Customer.Step
		resp.CustomerData = CustomerData{
			Email:    preMembers[0].Email,
			Mobile:   preMembers[0].Mobile,
			MemberID: preMembers[0].MemberID,
			Fullname: &CustomerFullname{
				ID:           preMembers[0].Customer.ID,
				MemberID:     preMembers[0].Customer.MemberID,
				Citizenship:  preMembers[0].Customer.Citizenship,
				Title:        preMembers[0].Customer.Title,
				ThName:       preMembers[0].Customer.ThName,
				ThMiddleName: preMembers[0].Customer.ThMiddleName,
				ThSurname:    preMembers[0].Customer.ThSurname,
				EnName:       preMembers[0].Customer.EnName,
				EnMiddleName: preMembers[0].Customer.EnMiddleName,
				EnSurname:    preMembers[0].Customer.EnSurname,
				Mobile:       preMembers[0].Customer.Mobile,
				Email:        preMembers[0].Customer.Email,
				Agreement:    preMembers[0].Customer.Agreement,
				Step:         preMembers[0].Customer.Step,
				Type:         preMembers[0].Customer.Type,
			},
			IDCard: &IDCardDetail{
				DateOfBirth: preMembers[0].Customer.DateOfBirth,
				Status:      preMembers[0].Customer.Status,
				IDCard:      preMembers[0].Customer.IDCard,
				LaserCode:   preMembers[0].Customer.LaserCode,
				ExpireDate:  preMembers[0].Customer.ExpireDate,
			},
			SuiteTest:     &preMembers[0].SuiteTest,
			Document:      &preMembers[0].Document,
			Addresses:     preMembers[0].Addresses,
			SourceOfFund:  &preMembers[0].SourceOfFund,
			Occupation:    &preMembers[0].Occupation,
			Banks:         preMembers[0].Banks,
			KnowledgeTest: preMembers[0].Customer.KnowledgeTest,
		}

		if preMembers[0].SuiteTest.ID == "" {
			resp.CustomerData.SuiteTest = nil
		}
		if preMembers[0].Document.ID == "" {
			resp.CustomerData.Document = nil
		}
		if preMembers[0].SourceOfFund.ID == "" {
			resp.CustomerData.SourceOfFund = nil
		}
		if preMembers[0].Occupation.ID == "" {
			resp.CustomerData.Occupation = nil
		}
		if len(preMembers[0].Banks) == 0 {
			resp.CustomerData.Banks = nil
		}
		if len(preMembers[0].Addresses) == 0 {
			resp.CustomerData.Addresses = nil
		}
	}

	resp.Member.ID = memberId
	resp.Member.Email = email
	resp.Member.Mobile = mobile

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
	appmans, err := u.getAppman(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	if err = u.upsertPreMember(data.Email, data.Mobile, memberId, appmans, ctx, tx); err != nil {
		return PreCitizenshipResp{}, err
	}

	preMembers, err := u.getPreMemberJoinedCustomer(memberId, ctx, tx)
	if err != nil {
		return PreCitizenshipResp{}, err
	}

	step, resp := u.mapToPreMemberResponse(preMembers, data.Email, data.Mobile, memberId)
	if len(appmans) > 0 {
		if step == 0 {
			step = 50
		}
		resp.Step = step
		u.setCitizenship(&resp, memberId, ctx, tx)
	}

	return
}
