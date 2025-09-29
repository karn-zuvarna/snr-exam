package onboarding

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go -package=mock
type IExternal interface {
	GetTimeNow() time.Time
	GenUuid() string
	GetIP(c echo.Context) string
}
type IServiceRiskScore interface {
	LocationAndCustomerRisk(ctx context.Context, tx *gorm.DB, funds *SourceOfFundsDB) (err error)
}

type IService interface {
	VerifyToken(tokenStr string, keyData []byte) (PreCitizenToken, error)
}

type Irepo struct {
	PreMember        IPreMember
	IDCard           IIDcard
	Customer         ICustomer
	Suitetest        ISuitetest
	Answer           IAnswer
	Appman           IAppman
	Document         IDocument
	Address          IAddress
	SourceOfFunds    ISourceOfFunds
	Bank             IBank
	Occupation       IOccupation
	Liveness         ILiveness
	LivenessAtt      ILivenessAttributes
	Recognition      IRecognition
	Watchlist        IWatchList
	CddStatus        CddStatus
	RiskOccupation   IRiskOccupation
	RiskCountry      IRiskCountry
	RiskScore        IRiskScore
	Join             IJoin
	Cdd              ICdd
	MasterLocation   IMasterLocation
	MasterCountry    IMasterCountry
	MasterOccupation IMasterOccupation
}

type PreCitizenToken struct {
	Email    string `json:"email"`
	Mobile   string `json:"phone"`
	MemberID int    `json:"account_id"`
	// AccountID int    `json:"account_id"`
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type IPreMember interface {
	GetAll(ctx context.Context, tx *gorm.DB, preMember *[]PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
	CreateAll(ctx context.Context, tx *gorm.DB, preMember *[]PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
	UpdateAll(ctx context.Context, tx *gorm.DB, preMember *PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
	JoinCustomer(ctx context.Context, tx *gorm.DB, preMember *[]PreMemberDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
}
type IIDcard interface{}
type ICustomer interface {
	GetAll(ctx context.Context, tx *gorm.DB, customer *[]CustomerDetailsDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
}
type ISuitetest interface{}
type IAnswer interface{}
type IAppman interface {
	GetAll(ctx context.Context, tx *gorm.DB, appman *[]AppmanDB, scope ...func(*gorm.DB) *gorm.DB) (err error)
}
type IDocument interface{}
type IAddress interface{}
type ISourceOfFunds interface{}
type IBank interface{}
type IOccupation interface{}
type ILiveness interface{}
type ILivenessAttributes interface{}
type IRecognition interface{}
type IWatchList interface{}
type CddStatus interface{}
type IRiskOccupation interface{}
type IRiskCountry interface{}
type IRiskScore interface{}
type IJoin interface{}
type ICdd interface{}
type IMasterLocation interface{}
type IMasterCountry interface{}
type IMasterOccupation interface{}
