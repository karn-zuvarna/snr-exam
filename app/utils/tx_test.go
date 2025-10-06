package utils

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OnboardingUnitTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *OnboardingUnitTestSuite) Test_RollbackOnError_NoError() {
	tx := s.db.Begin()
	defer tx.Commit()

	isError := RollbackOnError(tx, nil)
	s.Equal(nil, isError)
}

func (s *OnboardingUnitTestSuite) Test_RollbackOnError_HasError() {
	tx := s.db.Begin()
	defer tx.Commit()

	isError := RollbackOnError(tx, tx.Error)
	s.NotNil(isError)
}
