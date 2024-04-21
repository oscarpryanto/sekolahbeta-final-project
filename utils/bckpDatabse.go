package utils

import (
	"final-project/bckp-database/config"
	"final-project/bckp-database/model"
	"time"
)

func InsertBckpDatabase(data model.BckpDatabase) (model.BckpDatabase, error) {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := data.Create(config.Mysql.DB)

	return data, err
}

func LatestBackup(data model.BckpDatabase) ([]model.BckpDatabase, error) {
	latestBackup, err := data.LatestBackup(config.Mysql.DB)
	if err != nil {
		panic(err)
	}

	return latestBackup, err
}

func BackupHistoryByName(databaseName string) ([]model.BckpDatabase, error) {
	database := model.BckpDatabase{
		DatabaseName: databaseName,
	}

	backupHistory, err := database.GetHistoryBackup(config.Mysql.DB)
	if err != nil {
		panic(err)
	}

	return backupHistory, err
}

func BackupDatabaseById(databaseId uint) (model.BckpDatabase, error) {
	database := model.BckpDatabase{
		Model: model.Model{
			ID: databaseId,
		},
	}

	databaseInfo, err := database.GetBackupById(config.Mysql.DB)
	if err != nil {
		panic(err)
	}

	return databaseInfo, err
}
