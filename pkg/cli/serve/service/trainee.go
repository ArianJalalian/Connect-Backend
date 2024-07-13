package serve

import (
	"errors"
	"fmt"
	"time"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"gorm.io/gorm"
)

var ProgramStatuses = map[string]string{
	"TrainerPending":  "TrainerPending",
	"TraineePending":  "TraineePending",
	"TrainerRejected": "TrainerRejected",
	"TraineeRejected": "TraineeRejected",
	"TraineeAccepted": "TraineeAccepted",
	"Confirmed":       "Confirmed",
}

func GetTraineeProfile(id uint) (models.Trainee, error) {
	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainee{}, nil
		}
		return models.Trainee{}, err
	}
	return trainee, nil
}

func GetTraineeById(id uint) (models.Trainee, error) {
	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainee{}, nil
		}
		return models.Trainee{}, err
	}
	return trainee, nil
}

func EditTraineeProfile(id uint64, trainee dto.UserEditTraineeOrTrainer) (models.Trainee, error) {
	user, _ := GetUserById(id)

	activeDays := models.ActiveDays{
		// Monday:    trainee.ActiveDays[0],
		// Tuesday:   trainee.ActiveDays[1],
		// Wednesday: trainee.ActiveDays[2],
		// Thursday:  trainee.ActiveDays[3],
		// Friday:    trainee.ActiveDays[4],
		// Saturday:  trainee.ActiveDays[5],
		// Sunday:    trainee.ActiveDays[6],
	}

	if err := database.DB.Create(&activeDays).Error; err != nil {
		return models.Trainee{}, err
	}

	newTrainee := models.Trainee{
		UserID:         user.ID,
		Height:         trainee.Height,
		Weight:         trainee.Weight,
		UserName:       trainee.UserName,
		Status:         trainee.Status,
		Contact:        trainee.Contact,
		Language:       trainee.Language,
		Country:        trainee.Country,
		Sports:         trainee.Sports,
		MedicalHistory: trainee.MedicalHistory,
		ActiveDaysID:   activeDays.ID,
		//TrainerID:      0,
	}

	user.InfoID = fmt.Sprintf("%d", newTrainee.ID)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingTrainee models.Trainee
	result := tx.Where(models.Trainee{UserID: newTrainee.UserID}).First(&existingTrainee)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := tx.Create(&newTrainee).Error; err != nil {
			tx.Rollback()
			return models.Trainee{}, err
		}
	} else if result.Error == nil {
		if err := tx.Model(&existingTrainee).Updates(&newTrainee).Error; err != nil {
			tx.Rollback()
			return models.Trainee{}, err
		}
		newTrainee = existingTrainee
	} else {
		tx.Rollback()
		return models.Trainee{}, result.Error
	}

	user.InfoID = fmt.Sprintf("%d", newTrainee.ID)
	if err := tx.Model(&user).Updates(models.User{InfoID: user.InfoID}).Error; err != nil {
		tx.Rollback()
		return models.Trainee{}, err
	}

	tx.Commit()

	createdTrainee, err := GetTraineeById(newTrainee.ID)
	if err != nil {
		return models.Trainee{}, err
	}
	return createdTrainee, nil
}

func CreateProgramRequest(request dto.ProgramRequest) (models.Request, error) {
	trainee, err := GetTraineeById(request.TraineeID)
	if err != nil {
		return models.Request{}, err
	}
	activeDays := models.ActiveDays{
		Monday:    request.ActiveDays[0],
		Tuesday:   request.ActiveDays[1],
		Wednesday: request.ActiveDays[2],
		Thursday:  request.ActiveDays[3],
		Friday:    request.ActiveDays[4],
		Saturday:  request.ActiveDays[5],
		Sunday:    request.ActiveDays[6],
	}

	if err := database.DB.Create(&activeDays).Error; err != nil {
		return models.Request{}, err
	}
	newRequest := models.Request{
		TraineeID:    request.TraineeID,
		TrainerID:    request.TrainerID,
		Description:  request.Description,
		Date:         time.Now(),
		TraineeName:  trainee.User.FirstName + " " + trainee.User.LastName,
		Status:       "TrainerPending",
		Price:        0,
		ActiveDaysID: activeDays.ID,
		ActiveDays:   activeDays,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&newRequest).Error; err != nil {
		return models.Request{}, err
	}
	trainee.RequestID = newRequest.ID
	if err := tx.Model(&trainee).Updates(models.Trainee{RequestID: trainee.RequestID}).Error; err != nil {
		tx.Rollback()
		return models.Request{}, err
	}

	trainer, err := GetTrainerById(request.TrainerID)
	if err != nil {
		return models.Request{}, err
	}
	trainer.RequestIDs = append(trainer.RequestIDs, int64(newRequest.ID))
	if err := tx.Save(&trainer).Error; err != nil {
		tx.Rollback()
		return models.Request{}, err
	}

	tx.Commit()

	return GetRequest(newRequest.ID)
}

func GetRequest(id uint) (models.Request, error) {
	var req models.Request
	if err := database.DB.Preload("ActiveDays").Where("id = ?", id).First(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Request{}, nil
		}
		return models.Request{}, err
	}
	return req, nil
}

func ChangeStatus(change dto.TraineeChangeStatus) (models.Request, error) {
	if ProgramStatuses[change.Status] == "" {
		return models.Request{}, errors.New("invalid status")
	}
	r, err := GetRequest(change.RequestID)
	if err != nil {
		return models.Request{}, err
	}
	tx := database.DB.Begin()
	r.Status = change.Status
	if err := tx.Save(&r).Error; err != nil {
		tx.Rollback()
		return models.Request{}, err
	}
	tx.Commit()
	return GetRequest(r.ID)
}

func GetTraineeByUserID(id uint) (models.Trainee, error) {
	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("user_id = ?", id).First(&trainee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainee{}, nil
		}
		return models.Trainee{}, err
	}
	return trainee, nil
}

func GetTrainingProgramByRequestID(id uint) (*models.TrainingProgram, error) {
	var p models.TrainingProgram
	req, err := GetRequest(id)
	if err != nil {
		return nil, err
	}
	trainerID := req.TrainerID
	if err := database.DB.Where("trainer_id = ?", trainerID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
