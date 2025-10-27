package test


import (
	"go-crud/initializers"
	"go-crud/router"
	"testing"

	"github.com/gin-gonic/gin"
)


type BaseTestSuite struct {
	router *gin.Engine
	t *testing.T
}

func NewTestSuite(t *testing.T) *BaseTestSuite {
	suite := &BaseTestSuite{t: t}
	suite.SetUp()
	return suite
}

func (suite *BaseTestSuite) SetUp() {
	gin.SetMode(gin.TestMode)
	suite.router = router.SetupRouter()
	suite.CleanUp()
}

func (suite *BaseTestSuite) CleanUp() {
	// Force delete all records in correct order to avoid FK constraints
	// Use raw SQL to ensure complete cleanup regardless of relationships
	initializers.DB.Exec("DELETE FROM post_versions WHERE 1=1")
	initializers.DB.Exec("DELETE FROM posts WHERE 1=1")
	initializers.DB.Exec("DELETE FROM users WHERE 1=1")
}

func (suite *BaseTestSuite) TearDown() {
	suite.CleanUp()
}
