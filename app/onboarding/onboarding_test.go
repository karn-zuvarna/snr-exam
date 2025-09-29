package onboarding_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/karn-zuvarna/snr-exam/app/onboarding"
	"github.com/karn-zuvarna/snr-exam/app/onboarding/mock"
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
	repo    onboarding.Irepo

	uc *onboarding.Onboarding
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
	s.repo = onboarding.Irepo{
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
	s.uc = onboarding.New(s.ext, s.db, s.service, s.repo)
}

func (s *OnboardingTestSuite) TearDownTest() {

	s.ctrl.Finish()
}

func TestIntegratedTest(t *testing.T) {
	suite.Run(t, new(OnboardingTestSuite))
}

func (s *OnboardingTestSuite) Test_PreCitizenship__Success() {

	timeNow, _ := time.Parse("20060102", "20250129")
	zeroTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	s.ext.EXPECT().GetTimeNow().Return(timeNow)

	s.service.EXPECT().VerifyToken("token", []byte("public_key")).Return(onboarding.PreCitizenToken{
		Email:    "john.doe@example.com",
		Mobile:   "0888888888",
		MemberID: 1,
		UserID:   1,
	}, nil)

	s.repo.Appman.(*mock.MockIAppman).EXPECT().GetAll(s.ctx, s.db, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]onboarding.AppmanDB, scope ...func(*gorm.DB) *gorm.DB) error {
			*data = []onboarding.AppmanDB{
				{
					ID:             "appman-id-1",
					MemberID:       "1",
					VerificationID: "verification-1",
					UserID:         "1",
					IsFrontIDCard:  true,
					IsBackIDCard:   true,
					IsLiveness:     true,
					IsRecognition:  true,
					Status:         "completed",
					VerifiedStatus: true,
					DopaStatus:     true,
					FaceSimilarity: 0.95,
					Type:           1,
				},
				{
					ID:             "appman-id-2",
					MemberID:       "1",
					VerificationID: "verification-2",
					UserID:         "1",
					IsFrontIDCard:  true,
					IsBackIDCard:   true,
					IsLiveness:     true,
					IsRecognition:  true,
					Status:         "completed",
					VerifiedStatus: true,
					DopaStatus:     true,
					FaceSimilarity: 0.95,
					Type:           1,
				},
			}
			return nil
		})

	s.repo.PreMember.(*mock.MockIPreMember).EXPECT().JoinCustomer(s.ctx, s.db, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]onboarding.PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) error {
			*data = []onboarding.PreMemberDB{
				{
					BaseGorm: onboarding.BaseGorm{
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						DeletedAt: gorm.DeletedAt{},
					},
					ID:       "premember-001",
					Email:    "john.doe@example.com",
					Mobile:   "0888888888",
					MemberID: "1",

					Customer: onboarding.CustomerDetailsDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:                   "1",
						CreatedBy:            "system",
						MemberID:             "1",
						Citizenship:          1,
						Title:                "นาย",
						ThName:               "จอห์น",
						ThMiddleName:         "",
						ThSurname:            "โด",
						EnName:               "John",
						EnMiddleName:         "",
						EnSurname:            "Doe",
						Mobile:               "0888888888",
						Email:                "john.doe@example.com",
						Agreement:            true,
						DateOfBirth:          timeNow,
						Status:               1,
						IDCard:               "1234567890123",
						LaserCode:            "ABC123456789",
						ExpireDate:           timeNow,
						SourceOfFunds:        1,
						KnowledgeTest:        true,
						KnowledgeTestVersion: "v1.0",
						Step:                 5,
						Type:                 1,
						Identify:             1,
					},

					Appman: onboarding.AppmanDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:             "appman-id-1",
						MemberID:       "1",
						VerificationID: "verification-1",
						UserID:         "1",
						IsFrontIDCard:  true,
						IsBackIDCard:   true,
						IsLiveness:     true,
						IsRecognition:  true,
						Status:         "completed",
						VerifiedStatus: true,
						VerifiedAt:     &timeNow,
						ExpiresAt:      &timeNow,
						DopaMessage:    "Verification successful",
						DopaStatus:     true,
						FaceSimilarity: 0.95,
						Type:           1,
					},

					SuiteTest: onboarding.SuiteTestDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:          "suitetest-001",
						CreatedBy:   "system",
						MemberID:    "member-001",
						TotalScore:  85,
						Acknowledge: true,
						Answers: []onboarding.SuiteTestAnswersDB{
							{
								BaseGorm: onboarding.BaseGorm{
									CreatedAt: timeNow,
									UpdatedAt: timeNow,
									DeletedAt: gorm.DeletedAt{},
								},
								ID:          "answer-001",
								CreatedBy:   "system",
								ReferenceID: "suitetest-001",
								Question:    "What is the recommended diversification strategy?",
								Answer:      "Portfolio diversification across asset classes",
							},
							{
								BaseGorm: onboarding.BaseGorm{
									CreatedAt: timeNow,
									UpdatedAt: timeNow,
									DeletedAt: gorm.DeletedAt{},
								},
								ID:          "answer-002",
								CreatedBy:   "system",
								ReferenceID: "suitetest-001",
								Question:    "What is risk tolerance?",
								Answer:      "Ability to withstand investment losses",
							},
						},
					},

					Document: onboarding.DocumentsDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:        "document-001",
						CreatedBy: "system",
						MemberID:  "member-001",
						DocTypes:  "ID_CARD",
						FileName:  "id_card_front.jpg",
						FilePath:  "/documents/member-001/id_card_front.jpg",
						FileTypes: "image/jpeg",
					},

					Addresses: []onboarding.AddressesDB{
						{
							BaseGorm: onboarding.BaseGorm{
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
								DeletedAt: gorm.DeletedAt{},
							},
							ID:           "address-001",
							CreatedBy:    "system",
							MemberID:     "member-001",
							LocationCode: "10110",
							ProvinceCode: "10",
							Number:       "123",
							Floor:        "5",
							Village:      "Happy Village",
							Alley:        "Soi 15",
							Road:         "Sukhumvit Road",
							Subdistrict:  "Khlong Toei",
							District:     "Khlong Toei",
							Province:     "Bangkok",
							Zipcode:      "10110",
							CountryCode:  "TH",
							Type:         1,
						},
						{
							BaseGorm: onboarding.BaseGorm{
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
								DeletedAt: gorm.DeletedAt{},
							},
							ID:           "address-002",
							CreatedBy:    "system",
							MemberID:     "member-001",
							LocationCode: "10200",
							ProvinceCode: "10",
							Number:       "456",
							Floor:        "10",
							Village:      "Work Plaza",
							Alley:        "Soi 20",
							Road:         "Silom Road",
							Subdistrict:  "Silom",
							District:     "Bang Rak",
							Province:     "Bangkok",
							Zipcode:      "10200",
							CountryCode:  "TH",
							Type:         3,
						},
					},

					SourceOfFund: onboarding.SourceOfFundsDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:                         "sourcefund-001",
						CreatedBy:                  "system",
						MemberID:                   "member-001",
						SourceOfFunds:              "salary,investment",
						CountrySourceOfIncome:      "TH",
						PurposeOfInvestment:        "long_term_growth",
						CountrySourceIncomeVersion: "v1.0",
					},

					Occupation: onboarding.OccupationDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:                  "occupation-001",
						CreateBy:            "system",
						MemberID:            "member-001",
						Education:           "bachelor_degree",
						EducationVersion:    "v1.0",
						SourceOfInvestment:  "salary,savings",
						Occupation:          "software_engineer",
						OccupationName:      "Software Engineer",
						VersionOfOccupation: "v1.0",
						Workplace:           "Tech Company Ltd.",
						JobPosition:         "Senior Developer",
						BusinessType:        "technology",
						BusinessTypeID:      "BT001",
						VersionOfBusiness:   "v1.0",
						MonthlyIncome:       "50,001-100,000",
						WithdrawalLimit:     "1,000,000",
					},

					Banks: []onboarding.BankDB{
						{
							BaseGorm: onboarding.BaseGorm{
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
								DeletedAt: gorm.DeletedAt{},
							},
							ID:            "bank-001",
							CreatedBy:     "system",
							BankID:        1,
							MemberID:      "member-001",
							BankName:      "Kasikorn Bank",
							BranchName:    "Sukhumvit Branch",
							AccountNumber: "1234567890",
							BankVersion:   "v1.0",
							Type:          1,
						},
						{
							BaseGorm: onboarding.BaseGorm{
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
								DeletedAt: gorm.DeletedAt{},
							},
							ID:            "bank-002",
							CreatedBy:     "system",
							BankID:        2,
							MemberID:      "member-001",
							BankName:      "Bangkok Bank",
							BranchName:    "Silom Branch",
							AccountNumber: "0987654321",
							BankVersion:   "v1.0",
							Type:          2,
						},
					},

					RiskScore: onboarding.RiskScoreDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:                    "riskscore-001",
						MemberID:              "member-001",
						Occupation:            3,
						Politician:            0,
						Amlo:                  0,
						NotResidenceInThai:    0,
						CountryDetails:        "",
						ChangeToCash:          2,
						TransferAbleAsset:     1,
						GlobalProduct:         0,
						ProductHighRisk:       0,
						SouthernBorder:        0,
						SouthernBorderDetails: "",
						CountryRisk:           0,
						CountryRiskDetails:    "",
						ChannelRisk:           1,
						ChannelDetails:        "Online platform",
						SpecialHighRisk:       false,
						SpecialDetails:        "",
						Prohibited:            false,
						ProhibitDetails:       "",
					},

					Cdd: onboarding.CddDB{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:           "cdd-001",
						MemberID:     "member-001",
						DopaMessage:  "DOPA verification successful",
						DopaStatus:   true,
						DopaDatetime: timeNow,
						AmloMessage:  "AMLO check passed",
						AmloStatus:   true,
						AmloDatetime: timeNow,
						LedMessage:   "LED check completed",
						LedStatus:    true,
						LedDatetime:  timeNow,
						PepMessage:   "PEP check clear",
						PepStatus:    true,
						PepDatetime:  timeNow,
						MuleMessage:  "Mule account check passed",
						MuleStatus:   true,
						MuleDatetime: timeNow,
					},
				},
			}
			return nil
		},
	)

	s.repo.Customer.(*mock.MockICustomer).EXPECT().GetAll(s.ctx, s.db, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, db *gorm.DB, data *[]onboarding.CustomerDetailsDB, scope ...func(*gorm.DB) *gorm.DB) error {
			*data = []onboarding.CustomerDetailsDB{}
			return nil
		})

	// var actual, err = s.uc.PreCitizenship(s.ctx, "token", "public_key")
	var actual, err = s.uc.PreCitizenship(s.ctx, "token", "public_key")
	s.NoError(err)
	s.Equal(onboarding.PreCitizenshipResp{
		BaseGorm: onboarding.BaseGorm{
			CreatedAt: zeroTime,
			UpdatedAt: zeroTime,
			DeletedAt: gorm.DeletedAt{},
		},
		RequestID: "",
		Status:    "",
		Member: onboarding.Member{
			ID:          "1",
			Mobile:      "0888888888",
			Email:       "john.doe@example.com",
			Citizenship: 1,
		},
		Step: 5,
		CustomerData: onboarding.CustomerData{
			ID:       "",
			Email:    "john.doe@example.com",
			Mobile:   "0888888888",
			Column:   "",
			MemberID: "1",
			Fullname: &onboarding.CustomerFullname{
				BaseGorm: onboarding.BaseGorm{
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					DeletedAt: gorm.DeletedAt{},
				},
				ID:           "1",
				MemberID:     "1",
				Citizenship:  1,
				Title:        "นาย",
				ThName:       "จอห์น",
				ThMiddleName: "",
				ThSurname:    "โด",
				EnName:       "John",
				EnMiddleName: "",
				EnSurname:    "Doe",
				Mobile:       "0888888888",
				Email:        "john.doe@example.com",
				Agreement:    true,
				Step:         5,
				Type:         1,
			},
			IDCard: &onboarding.IDCardDetail{
				DateOfBirth: timeNow,
				Status:      1,
				IDCard:      "1234567890123",
				LaserCode:   "ABC123456789",
				ExpireDate:  timeNow,
			},
			SuiteTest: &onboarding.SuiteTestDB{
				BaseGorm: onboarding.BaseGorm{
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					DeletedAt: gorm.DeletedAt{},
				},
				ID:          "suitetest-001",
				CreatedBy:   "system",
				MemberID:    "member-001",
				TotalScore:  85,
				Acknowledge: true,
				Answers: []onboarding.SuiteTestAnswersDB{
					{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:          "answer-001",
						CreatedBy:   "system",
						DeletedBy:   "",
						ReferenceID: "suitetest-001",
						Question:    "What is the recommended diversification strategy?",
						Answer:      "Portfolio diversification across asset classes",
					},
					{
						BaseGorm: onboarding.BaseGorm{
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
							DeletedAt: gorm.DeletedAt{},
						},
						ID:          "answer-002",
						CreatedBy:   "system",
						DeletedBy:   "",
						ReferenceID: "suitetest-001",
						Question:    "What is risk tolerance?",
						Answer:      "Ability to withstand investment losses",
					},
				},
			},
			Document: &onboarding.DocumentsDB{
				BaseGorm: onboarding.BaseGorm{
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					DeletedAt: gorm.DeletedAt{},
				},
				ID:        "document-001",
				CreatedBy: "system",
				MemberID:  "member-001",
				DocTypes:  "ID_CARD",
				FileName:  "id_card_front.jpg",
				FilePath:  "/documents/member-001/id_card_front.jpg",
				FileTypes: "image/jpeg",
			},
			Addresses: []onboarding.AddressesDB{
				{
					BaseGorm: onboarding.BaseGorm{
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						DeletedAt: gorm.DeletedAt{},
					},
					ID:           "address-001",
					CreatedBy:    "system",
					DeletedBy:    "",
					MemberID:     "member-001",
					LocationCode: "10110",
					ProvinceCode: "10",
					Number:       "123",
					Floor:        "5",
					Village:      "Happy Village",
					Alley:        "Soi 15",
					Road:         "Sukhumvit Road",
					Subdistrict:  "Khlong Toei",
					District:     "Khlong Toei",
					Province:     "Bangkok",
					Zipcode:      "10110",
					CountryCode:  "TH",
					Type:         1,
				},
				{
					BaseGorm: onboarding.BaseGorm{
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						DeletedAt: gorm.DeletedAt{},
					},
					ID:           "address-002",
					CreatedBy:    "system",
					DeletedBy:    "",
					MemberID:     "member-001",
					LocationCode: "10200",
					ProvinceCode: "10",
					Number:       "456",
					Floor:        "10",
					Village:      "Work Plaza",
					Alley:        "Soi 20",
					Road:         "Silom Road",
					Subdistrict:  "Silom",
					District:     "Bang Rak",
					Province:     "Bangkok",
					Zipcode:      "10200",
					CountryCode:  "TH",
					Type:         3,
				},
			},
			SourceOfFund: &onboarding.SourceOfFundsDB{
				BaseGorm: onboarding.BaseGorm{
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					DeletedAt: gorm.DeletedAt{},
				},
				ID:                         "sourcefund-001",
				CreatedBy:                  "system",
				MemberID:                   "member-001",
				SourceOfFunds:              "salary,investment",
				CountrySourceOfIncome:      "TH",
				PurposeOfInvestment:        "long_term_growth",
				CountrySourceIncomeVersion: "v1.0",
			},
			Occupation: &onboarding.OccupationDB{
				BaseGorm: onboarding.BaseGorm{
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					DeletedAt: gorm.DeletedAt{},
				},
				ID:                  "occupation-001",
				CreateBy:            "system",
				MemberID:            "member-001",
				Education:           "bachelor_degree",
				EducationVersion:    "v1.0",
				SourceOfInvestment:  "salary,savings",
				Occupation:          "software_engineer",
				OccupationName:      "Software Engineer",
				VersionOfOccupation: "v1.0",
				Workplace:           "Tech Company Ltd.",
				JobPosition:         "Senior Developer",
				BusinessType:        "technology",
				BusinessTypeID:      "BT001",
				VersionOfBusiness:   "v1.0",
				MonthlyIncome:       "50,001-100,000",
				WithdrawalLimit:     "1,000,000",
			},
			Banks: []onboarding.BankDB{
				{
					BaseGorm: onboarding.BaseGorm{
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						DeletedAt: gorm.DeletedAt{},
					},
					ID:            "bank-001",
					CreatedBy:     "system",
					DeletedBy:     "",
					BankID:        1,
					MemberID:      "member-001",
					BankName:      "Kasikorn Bank",
					BranchName:    "Sukhumvit Branch",
					AccountNumber: "1234567890",
					BankVersion:   "v1.0",
					Type:          1,
				},
				{
					BaseGorm: onboarding.BaseGorm{
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
						DeletedAt: gorm.DeletedAt{},
					},
					ID:            "bank-002",
					CreatedBy:     "system",
					DeletedBy:     "",
					BankID:        2,
					MemberID:      "member-001",
					BankName:      "Bangkok Bank",
					BranchName:    "Silom Branch",
					AccountNumber: "0987654321",
					BankVersion:   "v1.0",
					Type:          2,
				},
			},
			KnowledgeTest: true,
		},
	}, actual)
}
