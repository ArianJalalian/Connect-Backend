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

func GetTrainerProfile(id uint) (models.Trainer, error) {
	var trainer models.Trainer
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainer{}, nil
		}
		return models.Trainer{}, err
	}
	return trainer, nil
}

func GetTrainerById(id uint) (models.Trainer, error) {
	var trainer models.Trainer
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainer{}, nil
		}
		return models.Trainer{}, err
	}
	return trainer, nil
}

func GetTrainerByUserID(id uint) (models.Trainer, error) {
	var trainer models.Trainer
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("user_id = ?", id).First(&trainer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainer{}, nil
		}
		return models.Trainer{}, err
	}
	return trainer, nil
}

func EditTrainerProfile(id uint64, trainer dto.UserEditTraineeOrTrainer) (models.Trainer, error) {
	user, _ := GetUserById(id)

	activeDays := models.ActiveDays{
		Monday:    trainer.ActiveDays[0],
		Tuesday:   trainer.ActiveDays[1],
		Wednesday: trainer.ActiveDays[2],
		Thursday:  trainer.ActiveDays[3],
		Friday:    trainer.ActiveDays[4],
		Saturday:  trainer.ActiveDays[5],
		Sunday:    trainer.ActiveDays[6],
	}
	if err := database.DB.Create(&activeDays).Error; err != nil {
		return models.Trainer{}, err
	}

	newTrainer := models.Trainer{
		UserID:          user.ID,
		UserName:        trainer.UserName,
		Status:          trainer.Status,
		CoachExperience: trainer.CoachExperience,
		Contact:         trainer.Contact,
		Language:        trainer.Language,
		Country:         trainer.Country,
		Sport:           trainer.Sport,
		Achievements:    trainer.Achievements,
		Education:       trainer.Education,
		ActiveDaysID:    activeDays.ID,
	}

	user.InfoID = fmt.Sprintf("%d", newTrainer.ID)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingTrainer models.Trainer
	result := tx.Where(models.Trainee{UserID: newTrainer.UserID}).First(&existingTrainer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := tx.Create(&newTrainer).Error; err != nil {
			tx.Rollback()
			return models.Trainer{}, err
		}
	} else if result.Error == nil {
		if err := tx.Model(&existingTrainer).Updates(&newTrainer).Error; err != nil {
			tx.Rollback()
			return models.Trainer{}, err
		}
		newTrainer = existingTrainer
	} else {
		tx.Rollback()
		return models.Trainer{}, result.Error
	}

	user.InfoID = fmt.Sprintf("%d", newTrainer.ID)
	if err := tx.Model(&user).Updates(models.User{InfoID: user.InfoID}).Error; err != nil {
		tx.Rollback()
		return models.Trainer{}, err
	}

	tx.Commit()

	createdTrainer, err := GetTrainerById(newTrainer.ID)
	if err != nil {
		return models.Trainer{}, err
	}
	return createdTrainer, nil
}

func GetRequests(trainer models.Trainer) ([]models.Request, error) {
	var requests []models.Request
	for _, i := range trainer.RequestIDs {
		r, err := GetRequest(uint(i))
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func SetPrice(setPrice dto.TrainerSetPrice) (models.Request, error) {
	r, err := GetRequest(setPrice.RequestId)
	if err != nil {
		return models.Request{}, err
	}
	tx := database.DB.Begin()
	if setPrice.Rejected {
		r.Status = ProgramStatuses["TrainerRejected"]
	} else {
		r.Price = setPrice.Price
		r.Status = ProgramStatuses["TraineePending"]
	}
	if err := tx.Save(&r).Error; err != nil {
		tx.Rollback()
		return models.Request{}, err
	}
	tx.Commit()
	return r, nil
}

func CreateTrainingProgram(program dto.TrainingProgram) (models.TrainingProgram, error) {
	request, _ := GetRequest(program.RequestID)
	if request.Status != ProgramStatuses["TraineeAccepted"] {
		return models.TrainingProgram{}, errors.New("invalid request")
	}
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, program.StartDate)
	if err != nil {
		return models.TrainingProgram{}, err
	}
	endDate, err := time.Parse(layout, program.EndDate)
	if err != nil {
		return models.TrainingProgram{}, err
	}
	newProgram := models.TrainingProgram{
		TraineeID:      request.TraineeID,
		TrainerID:      request.TrainerID,
		Title:          program.Title,
		Description:    program.Description,
		Price:          request.Price,
		StartDate:      startDate,
		EndDate:        endDate,
		ActivityDaysID: request.ActiveDays.ID,
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&newProgram).Error; err != nil {
		tx.Rollback()
		return models.TrainingProgram{}, err
	}
	request.Status = ProgramStatuses["Confirmed"]
	if err := tx.Save(&request).Error; err != nil {
		tx.Rollback()
		return models.TrainingProgram{}, err
	}

	tx.Commit()

	createdProgram, err := GetTrainingProgram(newProgram.ID)
	if err != nil {
		return models.TrainingProgram{}, err
	}
	return createdProgram, nil
}

func GetTrainingProgram(id uint) (models.TrainingProgram, error) {
	var p models.TrainingProgram
	if err := database.DB.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TrainingProgram{}, nil
		}
		return models.TrainingProgram{}, err
	}
	return p, nil
}

func AddSportActivity(activity dto.AddSportActivity) (models.SportActivity, error) {
	program, _ := GetTrainingProgram(activity.ProgramID)
	// var existingMedia models.Media
	// if err := database.DB.Where("id = ?", activity.SportActivit.Sport.VideoID).First(&existingMedia).Error; err != nil {
	// 	return models.TrainingProgram{}, fmt.Errorf("video does not exist")
	// }

	sport := models.Sport{
		Title:       activity.SportActivit.Sport.Title,
		Description: activity.SportActivit.Sport.Description,
		VideoID:     activity.SportActivit.Sport.VideoID,
	}

	if err := database.DB.Create(&sport).Error; err != nil {
		return models.SportActivity{}, err
	}
	sportActivity := models.SportActivity{
		OrderNumber:       activity.SportActivit.OrderNumber,
		ExpectedValue:     activity.SportActivit.ExpectedValue,
		Value:             activity.SportActivit.Value,
		Status:            "Not Done",
		TrainingProgramID: activity.ProgramID,
		Sport:             sport,
		SportID:           sport.ID,
	}
	program.SportActivitys = append(program.SportActivitys, sportActivity)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&sportActivity).Error; err != nil {
		tx.Rollback()
		return models.SportActivity{}, err
	}
	if err := tx.Save(&program).Error; err != nil {
		tx.Rollback()
		return models.SportActivity{}, err
	}

	tx.Commit()
	return sportActivity, nil
}

func GetALLTrainers() ([]models.Trainer, error) {
	var ts []models.Trainer
	err := database.DB.Preload("User").Find(&ts).Error
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func AddReport(report dto.Report) (*models.Report, error) {
	newReport := models.Report{
		Description: report.Description,
		UserID:      report.UserID,
	}
	if err := database.DB.Create(&newReport).Error; err != nil {
		return nil, err
	}
	return &newReport, nil
}
