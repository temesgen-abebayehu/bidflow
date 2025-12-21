package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	tokenManager *TokenManager
	secretKey    string
}

// SetupTest runs before each test in the suite
func (suite *AuthTestSuite) SetupTest() {
	suite.secretKey = "test_secret_key_12345"
	suite.tokenManager = NewTokenManager(suite.secretKey)
}

func (suite *AuthTestSuite) TestGenerateAndVerifyToken() {
	userID := "user-123"
	companyID := "company-456"
	role := "BIDDER"

	// 1. Generate Token
	token, err := suite.tokenManager.GenerateToken(userID, companyID, role)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	// 2. Verify Token
	claims, err := suite.tokenManager.VerifyToken(token)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)

	// 3. Check Claims
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), companyID, claims.CompanyID)
	assert.Equal(suite.T(), role, claims.Role)
}

func (suite *AuthTestSuite) TestExpiredToken() {
	userID := "user-123"
	companyID := "company-456"
	role := "BIDDER"
	
	// Create a token that expired 1 hour ago
	token, _ := suite.tokenManager.GenerateToken(userID, companyID, role)

	claims, err := suite.tokenManager.VerifyToken(token)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrExpiredToken, err)
	assert.Nil(suite.T(), claims)
}

func (suite *AuthTestSuite) TestInvalidToken() {
	invalidToken := "this.is.not.a.valid.token"

	claims, err := suite.tokenManager.VerifyToken(invalidToken)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrInvalidToken, err)
	assert.Nil(suite.T(), claims)
}

// Run the suite
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}