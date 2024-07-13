package dto

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Age         uint   `json:"age" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required,min=11"`
	InfoID      string `json:"info_id"`
	InfoType    string `json:"info_type" validate:"required"`
	Block       bool   `json:"block"`
	Wallet      uint64 `json:"wallet"`
}

type LogIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPasswordHash(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	return err
}

type SetProgram struct {
	Activity string `json:"activity,omitempty"`
	Count    uint   `json:"count,omitempty"`
	Weight   uint   `json:"weight,omitempty"`
	Time     uint   `json:"time,omitempty"`
}

type UserEdit struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Age         uint   `json:"age"`
	Gender      string `json:"gender"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	InfoID      string `json:"info_id"`
	InfoType    string `json:"info_type"`
	Block       bool   `json:"block"`
	Wallet      uint64 `json:"wallet"`
}

type TrainerEdit struct {
	User            UserEdit `json:"user"`
	UserName        string   `json:"user_name" validate:"required"`
	Status          string   `json:"status" validate:"required"`
	CoachExperience uint     `json:"coach_experience" validate:"required"`
	Contact         string   `json:"contact" validate:"required"`
	Language        string   `json:"language" validate:"required"`
	Country         string   `json:"country" validate:"required"`
	Sport           string   `json:"sport" validate:"required"`
	Achievements    string   `json:"achievements" validate:"required"`
	Education       string   `json:"education" validate:"required"`
	ActiveDays      []bool   `json:"active_days" validate:"required"`
}

type UserEditTraineeOrTrainer struct {
	User            UserEdit `json:"user"`
	UserName        string   `json:"user_name"`
	Status          string   `json:"status"`
	CoachExperience uint     `json:"coach_experience"`
	Contact         string   `json:"contact"`
	Language        string   `json:"language"`
	Country         string   `json:"country"`
	Sport           string   `json:"sport"`
	Achievements    string   `json:"achievements"`
	Education       string   `json:"education"`
	Height          uint     `json:"height"`
	Weight          uint     `json:"weight"`
	Sports          string   `json:"sports"`
	MedicalHistory  string   `json:"medicalHistory"`
	ActiveDays      []bool   `json:"active_days"`
}

type TraineeEdit struct {
	User           UserEdit `json:"user"`
	Height         uint     `json:"height" validate:"required"`
	Weight         uint     `json:"weight" validate:"required"`
	Sports         string   `json:"sports" validate:"required"`
	UserName       string   `json:"user_name" validate:"required"`
	Status         string   `json:"status" validate:"required"`
	Contact        string   `json:"contact" validate:"required"`
	Language       string   `json:"language" validate:"required"`
	Country        string   `json:"country" validate:"required"`
	MedicalHistory string   `json:"medicalHistory" validate:"required"`
	ActiveDays     []bool   `json:"active_days" validate:"required"`
}

type ProgramRequest struct {
	ID          uint   `json:"id"`
	TrainerID   uint   `json:"trainerID" validate:"required"`
	TraineeID   uint   `json:"traineeID"`
	Description string `json:"description"`
	ActiveDays  []bool `json:"active_days" validate:"required"`
}

type TrainerSetPrice struct {
	RequestId uint `json:"requestId" validate:"required"`
	Price     uint `json:"price" validate:"required"`
	Rejected  bool `json:"rejected"`
	TrainerID uint `json:"trainerID"`
}

type TrainingProgram struct {
	RequestID   uint   `json:"request_id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
	TrainerID   uint   `json:"trainer_id"`
}

type TraineeChangeStatus struct {
	RequestID uint   `json:"request_id"`
	Status    string `json:"status"`
}

type AddSportActivity struct {
	ProgramID    uint                `json:"program_id"`
	SportActivit SportActivitRequest `json:"sport_activity"`
}
type SportActivitRequest struct {
	OrderNumber   uint         `json:"order_number"`
	ExpectedValue uint         `json:"expected_value"`
	Value         uint         `json:"value"`
	Sport         SportRequest `json:"sport"`
}

type SportRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoID     uint   `json:"video_id"`
}

type Report struct {
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
}
