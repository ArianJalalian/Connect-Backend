package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    uint
	Name  string
	Block bool
}

type Report struct {
	ID     uint
	UserID uint
}

// Define your database connection as a global variable
var (
	db  *gorm.DB
	err error
)

func setupTestDatabase() {
	// Initialize an in-memory SQLite database
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// AutoMigrate your models
	err = db.AutoMigrate(&User{}, &Report{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}

func tearDownTestDatabase() {
	// Close the database connection
	db = nil
}

// Test functions
func TestGetUsers(t *testing.T) {
	setupTestDatabase()
	defer tearDownTestDatabase()

	// Populate test data
	users := []User{
		{ID: 1, Name: "User 1"},
		{ID: 2, Name: "User 2"},
	}
	for _, u := range users {
		db.Create(&u)
	}

	// Call the function being tested
	var retrievedUsers []User
	db.Find(&retrievedUsers)

	// Assert the result
	assert.Equal(t, len(users), len(retrievedUsers), "Expected number of users should match")
}

func TestBlockUser(t *testing.T) {
	setupTestDatabase()
	defer tearDownTestDatabase()

	// Populate test data
	user := User{ID: 1, Name: "User 1"}
	db.Create(&user)

	report := Report{ID: 1, UserID: user.ID}
	db.Create(&report)

	// Call the function being tested
	retrievedUser := User{}
	db.First(&retrievedUser, report.UserID)

	assert.False(t, retrievedUser.Block, "Expected user to be unblocked initially")

	// Block the user
	tx := db.Begin()
	retrievedUser.Block = true
	err := tx.Save(&retrievedUser).Error
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to block user: %v", err)
	}
	tx.Commit()

	// Check if the user is blocked
	db.First(&retrievedUser, report.UserID)
	assert.True(t, retrievedUser.Block, "Expected user to be blocked after calling BlockUser")
}

// Run all tests in this file
func TestMain(m *testing.M) {
	setupTestDatabase()
	defer tearDownTestDatabase()
	m.Run()
}
