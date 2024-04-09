package controllers

import (
	"final-project/bckp-database/model"
	"final-project/bckp-database/presenter"
	"final-project/bckp-database/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RouteBckpDatabase(app *fiber.App) {
	databaseGroup := app.Group("/bckp-database")
	databaseGroup.Get("/", GetAllLatestBackupDatabase)
	// databaseGroup.Get("/:db_name")
	databaseGroup.Post("/:db_name", UploadFileHandler)
	databaseGroup.Get("/:id_file/download", DownloadFileHandler)
}

func UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("zip_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			map[string]any{
				"message": "Invalid form file",
			},
		)
	}

	dbName := c.Params("db_name")

	fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

	bckpDatabase, errCreateBckpDatabse := utils.InsertBckpDatabase(model.BckpDatabase{
		DatabaseName: dbName,
		FileName:     file.Filename,
		FilePath:     fmt.Sprintf("./upload/%s", file.Filename),
	})

	if errCreateBckpDatabse != nil {
		logrus.Printf("Terjadi error : %s\n", errCreateBckpDatabse.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(map[string]any{
				"message": "Server Error",
			})
	}

	err = c.SaveFile(file, fmt.Sprintf("./upload/%d.zip", bckpDatabase.ID))

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.AddBckpDatabseSuccessResponse(&bckpDatabase))

}

func DownloadFileHandler(c *fiber.Ctx) error {

	id_file := c.Params("id_file")

	zipFilePath := fmt.Sprintf("./upload/%s.zip", id_file)

	// Set appropriate headers for download
	c.Set("Content-Disposition", "attachment; filename=example.zip")
	c.Set("Content-Type", "application/zip")

	fmt.Println(zipFilePath)
	// Send the ZIP file to the client
	err := c.SendFile(zipFilePath)
	if err != nil {
		// Handle any errors (e.g., file not found)
		c.Status(fiber.StatusNotFound).
			JSON(map[string]any{
				"message": fmt.Sprintf("Backup database dengan id = %s tidak ditemukan!", id_file),
			})
	}
	return nil
}

func GetAllLatestBackupDatabase(c *fiber.Ctx) error {

	bckpDatabase, errCreateBckpDatabse := utils.LatestBackup(model.BckpDatabase{})

	if errCreateBckpDatabse != nil {
		logrus.Printf("Terjadi error : %s\n", errCreateBckpDatabse.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(map[string]any{
				"message": fmt.Sprintf("Server Error = %s", errCreateBckpDatabse.Error()),
			})
	}

	var backupData []map[string]interface{}
	for _, latestBackup := range bckpDatabase {
		backupData = append(backupData, map[string]interface{}{
			"database_name": latestBackup.DatabaseName,
			"latest_backup": map[string]interface{}{
				"id":        latestBackup.ID,
				"file_name": latestBackup.FileName,
				"timestamp": latestBackup.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"data":    backupData,
		"message": "Success",
	})

}
