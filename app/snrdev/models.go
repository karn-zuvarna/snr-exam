package snrdev

import (
	"time"

	"gorm.io/gorm"
)

type PreCitizenshipResp struct {
	RequestID    string       `json:"request_id"`
	Status       string       `json:"status"`
	Member       Member       `json:"member"`
	Step         int          `json:"step"`
	CustomerData CustomerData `json:"customer_data"`
	BaseGorm
}

type BaseGorm struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Member struct {
	ID          string `json:"id"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Citizenship int    `json:"citizenship"`
}

type CustomerData struct {
	ID            string            `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	Email         string            `json:"email" gorm:"uniqueIndex:idx_member_email_mobile"`
	Mobile        string            `json:"mobile" gorm:"uniqueIndex:idx_member_email_mobile"`
	Column        string            `json:"column"`
	MemberID      string            `json:"member_id" gorm:"unique"`
	Fullname      *CustomerFullname `json:"fullname"`
	IDCard        *IDCardDetail     `json:"id_card"`
	SuiteTest     *SuiteTestDB      `json:"suite_test" gorm:"foreignKey:MemberID;references:MemberID"`
	Document      *DocumentsDB      `json:"document" gorm:"foreignKey:MemberID;references:MemberID"`
	Addresses     []AddressesDB     `json:"addresses" gorm:"foreignKey:MemberID;references:MemberID"`
	SourceOfFund  *SourceOfFundsDB  `json:"source_of_fund" gorm:"foreignKey:MemberID;references:MemberID"`
	Occupation    *OccupationDB     `json:"occupation" gorm:"foreignKey:MemberID;references:MemberID"`
	Banks         []BankDB          `json:"banks" gorm:"foreignKey:MemberID;references:MemberID"`
	KnowledgeTest bool              `json:"knowledge_test"`
}

type CustomerFullname struct {
	ID           string `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	MemberID     string `json:"member_id" gorm:"unique"`
	Citizenship  int    `json:"citizenship"`
	Title        string `json:"title"`
	ThName       string `json:"th_name"`
	ThMiddleName string `json:"th_middle_name"`
	ThSurname    string `json:"th_surname"`
	EnName       string `json:"en_name"`
	EnMiddleName string `json:"en_middle_name"`
	EnSurname    string `json:"en_surname"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email"`
	Agreement    bool   `json:"agreement"`
	Step         int    `json:"step"`
	Type         int    `json:"type"`
	BaseGorm
}

type IDCardDetail struct {
	DateOfBirth time.Time `json:"date_of_birth"`
	Status      int       `json:"status"`
	IDCard      string    `json:"id_card"`
	LaserCode   string    `json:"laser_code"`
	ExpireDate  time.Time `json:"expire_date"`
}

type SuiteTestDB struct {
	ID          string               `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	CreatedBy   string               `json:"create_by,omitempty" gorm:"column:created_by"`
	DeletedBy   string               `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	MemberID    string               `json:"member_id" gorm:"unique"`
	TotalScore  int                  `json:"total_score"`
	Acknowledge bool                 `json:"acknowledge"`
	Answers     []SuiteTestAnswersDB `gorm:"foreignKey:ReferenceID;references:ID" json:"answers"`
	BaseGorm
}

type SuiteTestAnswersDB struct {
	BaseGorm
	ID string `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	// MemberID    string `json:"memberId" gorm:"unique"`
	CreatedBy   string `json:"create_by,omitempty" gorm:"column:created_by"`
	DeletedBy   string `json:"deletedBy,omitempty" gorm:"column:deleted_by"`
	ReferenceID string `json:"referenceId"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
}

type DocumentsDB struct {
	ID        string `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	CreatedBy string `json:"create_by,omitempty" gorm:"column:created_by"`
	DeletedBy string `json:"deleted_by,omitempty" gorm:"column:deleted_by"`
	MemberID  string `json:"member_id"`
	DocTypes  string `json:"doc_types"`
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	FileTypes string `json:"file_types"`
	BaseGorm
}

type AddressesDB struct {
	ID           string `json:"id"`
	CreatedBy    string `json:"create_by"`
	DeletedBy    string `json:"deleted_by"`
	MemberID     string `json:"member_id" gorm:"index:idx_member_type,unique"`
	LocationCode string `json:"location_code"`
	ProvinceCode string `json:"province_code"`
	Number       string `json:"number"`
	Floor        string `json:"floor"`
	Village      string `json:"village"`
	Alley        string `json:"alley"`
	Road         string `json:"road"`
	Subdistrict  string `json:"subdistrict"`
	District     string `json:"district"`
	Province     string `json:"province"`
	Zipcode      string `json:"zipcode"`
	CountryCode  string `json:"country_code"`
	Type         int    `json:"type" gorm:"index:idx_member_type,unique"`
	BaseGorm
}

type SourceOfFundsDB struct {
	ID                         string `json:"id"`
	CreatedBy                  string `json:"create_by"`
	DeletedBy                  string `json:"deleted_by"`
	MemberID                   string `json:"member_id" gorm:"unique"`
	SourceOfFunds              string `json:"source_of_funds"`
	CountrySourceOfIncome      string `json:"country_source"`
	PurposeOfInvestment        string `json:"purpose_of_investment"`
	CountrySourceIncomeVersion string `json:"country_source_version"`
	BaseGorm
}

type OccupationDB struct {
	ID                  string
	CreateBy            string `json:"create_by"`
	DeletedBy           string `json:"deleted_by"`
	MemberID            string `json:"member_id" validate:"required"`
	Education           string `json:"education" validate:"required"`
	EducationVersion    string `json:"education_version" validate:"required"`
	SourceOfInvestment  string `json:"source_of_investment" validate:"required"`
	Occupation          string `json:"occupation" validate:"required"`
	OccupationName      string `json:"occupation_name" validate:"required"`
	VersionOfOccupation string `json:"version_of_occupation" validate:"required"`
	Workplace           string `json:"workplace" validate:"required"`
	JobPosition         string `json:"job_position" validate:"required"`
	BusinessType        string `json:"business_type" validate:"required"`
	BusinessTypeID      string `json:"business_type_id" validate:"required"`
	VersionOfBusiness   string `json:"version_of_business" validate:"required"`
	MonthlyIncome       string `json:"monthly_income" validate:"required"`
	WithdrawalLimit     string `json:"withdrawal_limit" validate:"required"`
	BaseGorm
}
type BankDB struct {
	ID            string `json:"id"`
	CreatedBy     string `json:"created_by"`
	DeletedBy     string `json:"deleted_by"`
	BankID        int    `json:"bank_id"`
	MemberID      string `json:"member_id" gorm:"uniqueIndex:idx_banks_member_id_type"`
	BankName      string `json:"bank_name"`
	BranchName    string `json:"branch_name"`
	AccountNumber string `json:"account_number"`
	BankVersion   string `json:"bank_version"`
	Type          int    `json:"type" gorm:"uniqueIndex:idx_banks_member_id_type"`
	BaseGorm
}

type AppmanDB struct {
	ID             string     `json:"id"`
	MemberID       string     `json:"memberId"`
	VerificationID string     `json:"verificationId"`
	UserID         string     `json:"userId"`
	IsFrontIDCard  bool       `json:"isFrontIdCard"`
	IsBackIDCard   bool       `json:"isBackIdCard"`
	IsLiveness     bool       `json:"isLiveness"`
	IsRecognition  bool       `json:"isRecognition"`
	Status         string     `json:"status"`
	VerifiedStatus bool       `json:"verifiedStatus"`
	VerifiedAt     *time.Time `json:"verifiedAt"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	DopaMessage    string     `json:"dopaMessage"`
	DopaStatus     bool       `json:"dopaStatus"`
	FaceSimilarity float64    `json:"faceSimilarity"`
	Type           int        `json:"type"`
	BaseGorm
}

type PreMemberDB struct {
	ID           string            `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	Email        string            `gorm:"uniqueIndex:idx_member_email_mobile"`
	Mobile       string            `gorm:"uniqueIndex:idx_member_email_mobile"`
	MemberID     string            `gorm:"unique"`
	Customer     CustomerDetailsDB `gorm:"foreignKey:MemberID;references:MemberID"`
	Appman       AppmanDB          `gorm:"foreignKey:MemberID;references:MemberID"`
	SuiteTest    SuiteTestDB       `gorm:"foreignKey:MemberID;references:MemberID"`
	Document     DocumentsDB       `gorm:"foreignKey:MemberID;references:MemberID"`
	Addresses    []AddressesDB     `gorm:"foreignKey:MemberID;references:MemberID"`
	SourceOfFund SourceOfFundsDB   `gorm:"foreignKey:MemberID;references:MemberID"`
	Occupation   OccupationDB      `gorm:"foreignKey:MemberID;references:MemberID"`
	Banks        []BankDB          `gorm:"foreignKey:MemberID;references:MemberID"`
	RiskScore    RiskScoreDB       `gorm:"foreignKey:MemberID;references:MemberID"`
	Cdd          CddDB             `gorm:"foreignKey:MemberID;references:MemberID"`
	BaseGorm
}

type CustomerDetailsDB struct {
	ID                   string    `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	CreatedBy            string    `json:"createBy,omitempty" gorm:"column:created_by"`
	DeletedBy            string    `json:"deletedBy,omitempty" gorm:"column:deleted_by"`
	MemberID             string    `json:"member_id" gorm:"uniqueIndex;index:idx_member_mobile_email,unique"`
	Citizenship          int       `json:"citizenship"`
	Title                string    `json:"title" `
	ThName               string    `json:"thName" `
	ThMiddleName         string    `json:"thMiddleName"`
	ThSurname            string    `json:"thSurname"`
	EnName               string    `json:"enName"`
	EnMiddleName         string    `json:"enMiddleName"`
	EnSurname            string    `json:"enSurname"`
	Mobile               string    `json:"mobile" gorm:"index:idx_member_mobile_email,unique"`
	Email                string    `json:"email" gorm:"index:idx_member_mobile_email,unique"`
	Agreement            bool      `json:"agreement"`
	DateOfBirth          time.Time `json:"dateOfBirth" `
	Status               int       `json:"status"`
	IDCard               string    `json:"idCard"`
	LaserCode            string    `json:"laserCode"`
	ExpireDate           time.Time `json:"expireDate"`
	SourceOfFunds        int       `json:"sourceOfFunds"`
	KnowledgeTest        bool      `json:"KnowledgeTest"`
	KnowledgeTestScore   int       `json:"KnowledgeTestScore"`
	KnowledgeTestVersion string    `json:"KnowledgeTestVersion"`
	Step                 int       `json:"step"`
	Type                 int       `json:"type"`
	Identify             int       `json:"identify"`
	IP                   string    `json:"ip"`
	BaseGorm
}

type RiskScoreDB struct {
	ID                 string `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	MemberID           string `json:"memberId" gorm:"unique"`
	Occupation         int    `json:"occupation" gorm:"column:occupation"`
	Politician         int    `json:"politician" gorm:"column:politician"`
	Amlo               int    `json:"amlo" gorm:"column:amlo"`
	NotResidenceInThai int    `json:"notResidenceInThai" gorm:"column:not_residence_in_thai"`
	Ccib               int    `json:"ccib" gorm:"column:ccib"`
	CountryDetails     string `json:"countryDetails" gorm:"column:country_details"`

	ChangeToCash      int `json:"changeToCash" gorm:"column:change_to_cash"`
	TransferAbleAsset int `json:"transferAbleAsset" gorm:"column:transfer_able_asset"`
	GlobalProduct     int `json:"globalProduct" gorm:"column:global_product"`
	ProductHighRisk   int `json:"productHighRisk" gorm:"column:product_high_risk"`

	SouthernBorder        int    `json:"southern_border" gorm:"column:southern_border"`
	SouthernBorderDetails string `json:"southern_border_details" gorm:"column:southern_border_details"`
	CountryRisk           int    `json:"country_risk" gorm:"column:country_risk"`
	CountryRiskDetails    string `json:"country_risk_details" gorm:"column:country_risk_details"`

	ChannelRisk    int    `json:"ChannelRisk" gorm:"column:channel_risk"`
	ChannelDetails string `json:"channelDetails" gorm:"column:channel_details"`

	SpecialHighRisk bool   `json:"specialRisk" gorm:"column:special_high_risk"`
	SpecialDetails  string `json:"specialDetails" gorm:"column:special_details"`
	Prohibited      bool   `json:"prohibited" gorm:"column:prohibited"`
	ProhibitDetails string `json:"prohibitDetails" gorm:"column:prohibit_details"`
	Comment         string `json:"comment" gorm:"column:comment"`
	RiskLevel       int    `json:"risk_level" gorm:"column:risk_level"`
	BaseGorm
}

type CddDB struct {
	ID           string    `json:"id,omitempty" gorm:"primaryKey;type:varchar(40);not null"`
	MemberID     string    `json:"memberId" gorm:"unique"`
	DopaMessage  string    `json:"dopaMessage" gorm:"column:dopa_message"`
	DopaStatus   bool      `json:"dopaStatus" gorm:"column:dopa_status"`
	DopaDatetime time.Time `json:"dopaDatetime,omitempty" gorm:"column:dopa_datetime"`
	AmloMessage  string    `json:"amloMessage" gorm:"column:amlo_message"`
	AmloStatus   bool      `json:"amloStatus" gorm:"column:amlo_status"`
	AmloDatetime time.Time `json:"amloDatetime,omitempty" gorm:"column:amlo_datetime"`
	LedMessage   string    `json:"ledMessage" gorm:"column:led_message"`
	LedStatus    bool      `json:"ledStatus" gorm:"column:led_status"`
	LedDatetime  time.Time `json:"ledDatetime,omitempty" gorm:"column:led_datetime"`
	PepMessage   string    `json:"pepMessage" gorm:"column:pep_message"`
	PepStatus    bool      `json:"pepStatus" gorm:"column:pep_status"`
	PepDatetime  time.Time `json:"pepDatetime,omitempty" gorm:"column:pep_datetime"`
	MuleMessage  string    `json:"muleMessage" gorm:"column:mule_message"`
	MuleStatus   bool      `json:"muleStatus" gorm:"column:mule_status"`
	MuleDatetime time.Time `json:"muleDatetime,omitempty" gorm:"column:mule_datetime"`
	BaseGorm
}
