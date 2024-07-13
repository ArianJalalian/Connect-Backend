package dto

import "time"

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
}

type UserResponse struct {
	ID          uint              `json:"ID,omitempty"`
	Email       string            `json:"email,omitempty"`
	FirstName   string            `json:"first_name,omitempty"`
	LastName    string            `json:"last_name,omitempty"`
	Password    string            `json:"password,omitempty"`
	PhoneNumber string            `json:"phone_number,omitempty"`
	Type        string            `json:"type,omitempty"`
	InfoID      string            `json:"info_id,omitempty"`
	InfoType    string            `json:"info_type,omitempty"`
	Block       bool              `json:"block,omitempty"`
	Wallet      uint64            `json:"wallet,omitempty"`
	JWT         string            `json:"jwt,omitempty"`
	Profile     map[string]string `json:"profile,omitempty"`
}

type TrainerResponse struct {
	TrainerProfileCard `json:"trainer_profile_card"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	Sports             string `json:"sports,omitempty"`
	Achievements       string `json:"achievements,omitempty"`
	Education          string `json:"education,omitempty"`
	ID                 uint   `json:"id,omitempty"`
}

type TrainerProfileCard struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	UserName        string `json:"username,omitempty"`
	Email           string `json:"email,omitempty"`
	Status          string `json:"status,omitempty"`
	CoachExperience uint   `json:"coach_experience,omitempty"`
	Contact         string `json:"contact,omitempty"`
	Language        string `json:"language,omitempty"`
	Country         string `json:"country,omitempty"`
}

type TraineeInTrainerPage struct {
	Name string `json:"name,omitempty"`
}

type RequestsInTrainerPage struct {
	TraineeName string    `json:"trainee_name,omitempty"`
	Date        time.Time `json:"date"`
	Price       uint      `json:"price,omitempty"`
	Status      string    `json:"status,omitempty"`
	TrainerName string    `json:"trainer_name,omitempty"`
	TrainerID   uint      `json:"trainer_id,omitempty"`
}

type TrainerPlan struct {
	Monday    bool `json:"monday,omitempty"`
	Tuesday   bool `json:"tuesday,omitempty"`
	Wednesday bool `json:"wednesday,omitempty"`
	Thursday  bool `json:"thursday,omitempty"`
	Friday    bool `json:"friday,omitempty"`
	Saturday  bool `json:"saturday,omitempty"`
	Sunday    bool `json:"sunday,omitempty"`
}

type SportActivity struct {
	ExpectedCount  uint   `json:"expected_count,omitempty"`
	ExpectedWeight uint   `json:"expected_weight,omitempty"`
	MyCount        uint   `json:"my_count,omitempty"`
	Video          string `json:"video,omitempty"` // url
}

type Sport struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoPath   string `json:"video_path"`
}

type TraineeProgram struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Payment   uint      `json:"payment,omitempty"`
	Status    string    `json:"status,omitempty"`
}

type TraineeResponse struct {
	TraineeProfileCard `json:"trainee_profile_card"`
	SportExperience    string `json:"sport_experience,omitempty"`
	HealthProblems     string `json:"health_problems,omitempty"`
}

type TraineeProfileCard struct {
	UserName string `json:"userName,omitempty"`
	Email    string `json:"email,omitempty"`
	Status   string `json:"status,omitempty"`
	Wallet   uint64 `json:"wallet,omitempty"`
	Contact  string `json:"contact,omitempty"`
	Language string `json:"language,omitempty"`
	Country  string `json:"country,omitempty"`
}

type ProgramRequestSetPrice struct {
	ID          uint   `json:"id"`
	TrainerID   uint   `json:"trainerID" validate:"required"`
	TraineeID   uint   `json:"traineeID" validate:"required"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Status      string `json:"status"`
}

type SportResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoID     uint   `json:"video_id"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type ReportResponse struct {
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
	ReportID    uint   `json:"report_id"`
}

type WeekPlan struct {
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
	Sunday    bool `json:"sunday"`
}
