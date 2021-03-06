package models

import (
	"github.com/jinzhu/gorm"
)

// ScreenshotService ...
type ScreenshotService struct {
	DB *gorm.DB
	UpdatableModelService
}

// BatchCreate ...
func (s *ScreenshotService) BatchCreate(screenshots []*Screenshot) ([]*Screenshot, []error, error) {
	tx := s.DB.Begin()
	for _, screenshot := range screenshots {
		result := tx.Create(screenshot)
		verrs := ValidationErrors(result.GetErrors())
		if len(verrs) > 0 {
			tx.Rollback()
			return nil, verrs, nil
		}
		if result.Error != nil {
			tx.Rollback()
			return nil, nil, result.Error
		}

		err := tx.Preload("AppVersion").Preload("AppVersion.App").First(screenshot).Error
		if err != nil {
			return nil, nil, err
		}
	}
	return screenshots, nil, tx.Commit().Error
}

// Find ...
func (s *ScreenshotService) Find(screenshot *Screenshot) (*Screenshot, error) {
	err := s.DB.Where(screenshot).First(screenshot).Error

	if err != nil {
		return nil, err
	}

	return screenshot, nil
}

// FindAll ...
func (s *ScreenshotService) FindAll(appVersion *AppVersion) ([]Screenshot, error) {
	var screenshots []Screenshot
	err := s.DB.Preload("AppVersion").Preload("AppVersion.App").Where(map[string]interface{}{"app_version_id": appVersion.ID}).
		Find(&screenshots).
		Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return screenshots, nil
}

// BatchUpdate ...
func (s *ScreenshotService) BatchUpdate(screenshots []Screenshot, whitelist []string) ([]error, error) {
	for _, screenshot := range screenshots {
		updateData, err := s.UpdateData(screenshot, whitelist)
		if err != nil {
			return nil, err
		}
		result := s.DB.Model(&screenshot).Updates(updateData)
		verrs := ValidationErrors(result.GetErrors())
		if len(verrs) > 0 {
			return verrs, nil
		}
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return nil, nil
}

// Delete ...
func (s *ScreenshotService) Delete(screenshot *Screenshot) error {
	result := s.DB.Delete(&screenshot)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
