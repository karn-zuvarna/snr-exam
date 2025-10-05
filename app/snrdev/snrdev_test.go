package snrdev_test

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/karn-zuvarna/snr-exam/app/snrdev"
	"github.com/karn-zuvarna/snr-exam/app/snrdev/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OnboardingTestSuite struct {
	suite.Suite
	ctx  context.Context
	ctrl *gomock.Controller

	ext     *mock.MockIExternal
	db      *gorm.DB
	service *mock.MockIService
	repo    snrdev.Irepo

	uc *snrdev.Onboarding
}

func setupMock() *gorm.DB {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		log.Printf("failed to create mock database: %v", err)
	}

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.23"))

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
	}

	return db
}

func (s *OnboardingTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.ctx = context.Background()
	s.db = setupMock()
	s.ext = mock.NewMockIExternal(s.ctrl)
	s.service = mock.NewMockIService(s.ctrl)
	s.repo = snrdev.Irepo{
		Customer:       mock.NewMockICustomer(s.ctrl),
		PreMember:      mock.NewMockIPreMember(s.ctrl),
		IDCard:         mock.NewMockIIDcard(s.ctrl),
		Suitetest:      mock.NewMockISuitetest(s.ctrl),
		Answer:         mock.NewMockIAnswer(s.ctrl),
		Appman:         mock.NewMockIAppman(s.ctrl),
		Document:       mock.NewMockIDocument(s.ctrl),
		Address:        mock.NewMockIAddress(s.ctrl),
		SourceOfFunds:  mock.NewMockISourceOfFunds(s.ctrl),
		Bank:           mock.NewMockIBank(s.ctrl),
		Occupation:     mock.NewMockIOccupation(s.ctrl),
		Liveness:       mock.NewMockILiveness(s.ctrl),
		LivenessAtt:    mock.NewMockILivenessAttributes(s.ctrl),
		Recognition:    mock.NewMockIRecognition(s.ctrl),
		MasterLocation: mock.NewMockIMasterLocation(s.ctrl),
		RiskScore:      mock.NewMockIRiskScore(s.ctrl),
		Cdd:            mock.NewMockICdd(s.ctrl),
	}
	s.uc = snrdev.New(s.ext, s.db, s.service, s.repo)
}

func (s *OnboardingTestSuite) TearDownTest() {

	s.ctrl.Finish()
}

func TestIntegratedTest(t *testing.T) {
	suite.Run(t, new(OnboardingTestSuite))
}

func (s *OnboardingTestSuite) Test_Snr_Dev_Exam_Success() {
	var scopeCustomerData []func(*gorm.DB) *gorm.DB
	var preMembers []snrdev.PreMemberDB

	s.service.EXPECT().VerifyToken("token", []byte("public_key")).Return(precitizenToken, nil).Times(1)

	scopeCustomerData = append(scopeCustomerData, snrdev.Where("member_id = ?", 1))
	s.repo.PreMember.(*mock.MockIPreMember).EXPECT().JoinCustomer(s.ctx, s.db, &preMembers, gomock.AssignableToTypeOf(scopeCustomerData)).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]snrdev.PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) error {
			s.Require().Len(scope, 1)
			*data = preCitizenshipPreMemberGetAllReturn
			return nil
		},
	)

	var customerDetails []snrdev.CustomerDetailsDB
	s.repo.Customer.(*mock.MockICustomer).EXPECT().GetAll(s.ctx, s.db, &customerDetails, gomock.Any()).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]snrdev.CustomerDetailsDB, scope ...func(*gorm.DB) *gorm.DB) error {
			*data = []snrdev.CustomerDetailsDB{}
			return nil
		})

	var actual, err = s.uc.Snr_Dev_Exam(s.ctx, "token", "public_key")
	s.NoError(err)
	s.Equal(preCitizenshipExpectedResponse, actual)
}

func (s *OnboardingTestSuite) Test_getScopeAppman_Success() {
	var scopeAppman, scopesAppman []func(*gorm.DB) *gorm.DB
	var appman []snrdev.AppmanDB
	tx := s.db.Begin()
	defer tx.Commit()

	scopesAppman = append(scopesAppman, snrdev.Where("member_id = ?", 1))
	scopesAppman = append(scopesAppman, snrdev.Where("dopa_status = ?", true))
	scopesAppman = append(scopesAppman, snrdev.Order("updated_at desc"))
	s.repo.Appman.(*mock.MockIAppman).EXPECT().GetAll(s.ctx, s.db, &appman, gomock.AssignableToTypeOf(scopesAppman)).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]snrdev.AppmanDB, scope ...func(*gorm.DB) *gorm.DB) error {
			s.Require().Len(scope, 3)
			*data = preCitizenshipAppmanGetAllReturn
			return nil
		})

	actual, err := s.uc.TestGetScopeAppman(scopeAppman, "1", s.ctx, s.db)
	s.NoError(err)
	s.Equal(preCitizenshipAppmanGetAllReturn, actual)
}

func (s *OnboardingTestSuite) Test_upsertPremember_Success() {
	var queryPrememberGetAll = []snrdev.PreMemberDB{}
	var scopesAppman []func(*gorm.DB) *gorm.DB
	var appman []snrdev.AppmanDB
	tx := s.db.Begin()
	defer tx.Commit()

	scopesAppman = append(scopesAppman, snrdev.Where("member_id = ?", 1))
	scopesAppman = append(scopesAppman, snrdev.Where("dopa_status = ?", true))
	scopesAppman = append(scopesAppman, snrdev.Order("updated_at desc"))
	s.ext.EXPECT().GenUuid().AnyTimes().Return(uuid).Times(1)
	s.repo.PreMember.(*mock.MockIPreMember).EXPECT().GetAll(s.ctx, s.db, &queryPrememberGetAll, gomock.AssignableToTypeOf(scopesAppman)).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]snrdev.PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) error {
			s.Require().Len(scope, 2)
			*data = preCitizenshipPreMemberGetAllReturn
			return nil
		})
	s.repo.PreMember.(*mock.MockIPreMember).EXPECT().UpdateAll(s.ctx, gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf(scopesAppman)).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *snrdev.PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) error {
			*data = preCitizenshipPreMemberGetAllReturn[0]
			return nil
		})
	err := s.uc.TestUpsertPremember(s.ext.GenUuid(), precitizenToken.Email, precitizenToken.Mobile, strconv.Itoa(precitizenToken.MemberID), appman, s.ctx, tx)
	s.NoError(err)
}
