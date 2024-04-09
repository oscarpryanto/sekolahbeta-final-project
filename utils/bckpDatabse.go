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

	// data := map[string]interface{}{
	// 	"database_name": latestBackup.DatabaseName,
	// 	"latest_backup": map[string]interface{}{
	// 		"id":        latestBackup.ID,
	// 		"file_name": latestBackup.FileName,
	// 		"timestamp": latestBackup.CreatedAt.Format("2006-01-02 15:04:05"),
	// 	},
	// }

	return latestBackup, err
}
