package model

import (
	"gorm.io/gorm"
)

type BckpDatabase struct {
	Model
	DatabaseName string `gorm:"not null" json:"database_name"`
	FileName     string `gorm:"not null" json:"file_name"`
	FilePath     string `gorm:"not null" json:"file_path"`
}

func (cr *BckpDatabase) Create(db *gorm.DB) error {
	err := db.
		Model(BckpDatabase{}).
		Create(&cr).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (b *BckpDatabase) LatestBackup(db *gorm.DB) ([]BckpDatabase, error) {
	var latestBackup []BckpDatabase

	if err := db.Model(BckpDatabase{}).Where("deleted_at IS NULL").Select("created_at,database_name,file_name,file_path, MAX(id) as id").Order("created_at DESC").Group("database_name").Find(&latestBackup).Error; err != nil {
		return nil, err
	}

	return latestBackup, nil
}

func (cr *BckpDatabase) GetHistoryBackup(db *gorm.DB) ([]BckpDatabase, error) {
	res := []BckpDatabase{}

	if err := db.Model(BckpDatabase{}).Where("database_name = ?", cr.DatabaseName).Order("created_at DESC").Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil

}

func (cr *BckpDatabase) GetBackupById(db *gorm.DB) (BckpDatabase, error) {
	res := BckpDatabase{}

	if err := db.Model(BckpDatabase{}).Where("id = ?", cr.Model.ID).Find(&res).Error; err != nil {
		return BckpDatabase{}, err
	}

	return res, nil

}
