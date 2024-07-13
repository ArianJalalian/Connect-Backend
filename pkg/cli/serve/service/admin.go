package serve

import (
	"errors"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"gorm.io/gorm"
)

func GetUsers() []models.User {
	var users []models.User
	database.DB.Find(&users)
	return users
}

func AddSport(s dto.Sport) (*models.Sport, error) {
	video := models.Media{
		Path:      s.VideoPath,
		Name:      s.Title,
		Type:      "video",
		Size:      0,
		MediaType: "video",
	}
	if err := database.DB.Create(&video).Error; err != nil {
		return nil, err
	}
	video_res, err := GetVideoByID(video.ID)
	if err != nil {
		return nil, errors.New("video does not exist")
	}
	sport := models.Sport{
		Title:       s.Title,
		Description: s.Description,
		VideoID:     video_res.ID,
		Video:       *video_res,
	}
	if err := database.DB.Create(&sport).Error; err != nil {
		return nil, err
	}
	return GetSportByID(sport.ID)
}

func GetVideoByID(id uint) (*models.Media, error) {
	var media models.Media
	if err := database.DB.Where("id = ?", id).First(&media).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &media, nil
}
func GetSportByID(id uint) (*models.Sport, error) {
	var sport models.Sport
	if err := database.DB.Preload("Video").Where("id = ?", id).First(&sport).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sport, nil
}

func GetReports() []models.Report {
	var reports []models.Report
	database.DB.Find(&reports)
	return reports
}

func BlockUser(reportID uint) (*models.Report, error) {
	var report models.Report
	if err := database.DB.Where("id = ?", reportID).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	tx := database.DB.Begin()
	user, _ := GetUserById(uint64(report.UserID))
	user.Block = true
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &report, nil
}

func GetSports() ([]models.Sport, error) {
	var sports []models.Sport
	database.DB.Find(&sports)
	return sports, nil
}
