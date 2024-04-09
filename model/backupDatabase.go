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

func (cr *BckpDatabase) GetAllBckUp(db *gorm.DB) ([]BckpDatabase, error) {
	res := []BckpDatabase{}

	err := db.
		Model(BckpDatabase{}).
		Find(&res).
		Error

	if err != nil {
		return []BckpDatabase{}, err
	}

	return res, nil
}

// func (cr *Car) GetByID(db *gorm.DB) (Car, error) {
// 	res := Car{}

// 	err := db.
// 		Model(Car{}).
// 		Where("id = ?", cr.Model.ID).
// 		Take(&res).
// 		Error

// 	if err != nil {
// 		return Car{}, err
// 	}

// 	return res, nil
// }

// func (cr *Car) GetAll(db *gorm.DB) ([]Car, error) {
// 	res := []Car{}

// 	err := db.
// 		Model(Car{}).
// 		Find(&res).
// 		Error

// 	if err != nil {
// 		return []Car{}, err
// 	}

// 	return res, nil
// }

// func (cr *Car) UpdateOneByID(db *gorm.DB) error {
// 	err := db.
// 		Model(Car{}).
// 		Select("nama", "tipe", "tahun").
// 		Where("id = ?", cr.Model.ID).
// 		Updates(map[string]any{
// 			"nama":  cr.Nama,
// 			"tipe":  cr.Tipe,
// 			"tahun": cr.Tahun,
// 		}).
// 		Error

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (cr *Car) DeleteByID(db *gorm.DB) error {
// 	err := db.
// 		Model(Car{}).
// 		Where("id = ?", cr.Model.ID).
// 		Delete(&cr).
// 		Error

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
