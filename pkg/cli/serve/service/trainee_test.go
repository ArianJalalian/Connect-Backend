package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
)

func setupTestDatabase2() {
	// Initialize an in-memory SQLite database
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// AutoMigrate your models
	err = db.AutoMigrate(
		&models.User{},
		&models.Trainee{},
		&models.Request{},
		&models.ActiveDays{},
		// Add other models as needed
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Set the database connection for the package
	database.DB = db
}

func tearDownTestDatabase2() {
	// Close the database connection
	db = nil
}

// Run all tests in this file
func TestMain2(t *testing.T) {
	setupTestDatabase2()
	defer tearDownTestDatabase2()
	t.Run("EditTraineeProfile", TestEditTraineeProfile)
	t.Run("GetTraineeById", TestGetTraineeById)
	t.Run("GetTraineeByUserID", TestGetTraineeByUserID)
}

func TestEditTraineeProfile(t *testing.T) {
	// Setup test database
	setupTestDatabase2()
	defer tearDownTestDatabase2()

	// Create a test user
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		// Add other fields as needed
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Prepare test data for trainee profile edit
	traineeData := dto.UserEditTraineeOrTrainer{
		Height: 180,
		Weight: 75,
		// Add other fields as needed
	}

	// Call the function being tested
	createdTrainee, err := EditTraineeProfile(uint64(user.ID), traineeData)
	if err != nil {
		t.Fatalf("Failed to edit trainee profile: %v", err)
	}

	// Assert the created trainee
	assert.Equal(t, traineeData.Height, createdTrainee.Height, "Expected height to match")
	assert.Equal(t, traineeData.Weight, createdTrainee.Weight, "Expected weight to match")
	// Add more assertions as needed
}

func TestGetTraineeById(t *testing.T) {
	// Setup test database
	setupTestDatabase()
	defer tearDownTestDatabase()

	// Create a test trainee
	trainee := models.Trainee{
		UserID: 1,
		// Add other fields as needed
	}
	if err := db.Create(&trainee).Error; err != nil {
		t.Fatalf("Failed to create trainee: %v", err)
	}

	// Call the function being tested
	retrievedTrainee, err := GetTraineeById(trainee.ID)
	if err != nil {
		t.Fatalf("Failed to get trainee by ID: %v", err)
	}

	// Assert the retrieved trainee
	assert.Equal(t, trainee.UserID, retrievedTrainee.UserID, "Expected user ID to match")
	// Add more assertions as needed
}

func TestGetTraineeByUserID(t *testing.T) {
	// Setup test database
	setupTestDatabase()
	defer tearDownTestDatabase()

	// Create a test user
	user := models.User{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		// Add other fields as needed
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a trainee associated with the test user
	trainee := models.Trainee{
		UserID: user.ID,
		// Add other fields as needed
	}
	if err := db.Create(&trainee).Error; err != nil {
		t.Fatalf("Failed to create trainee: %v", err)
	}

	// Call the function being tested
	retrievedTrainee, err := GetTraineeByUserID(user.ID)
	if err != nil {
		return
	}
	fmt.Println(retrievedTrainee)
}
