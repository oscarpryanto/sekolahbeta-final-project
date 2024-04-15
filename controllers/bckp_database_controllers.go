package controllers

import (
	"errors"
	"final-project/bckp-database/model"
	"final-project/bckp-database/presenter"
	"final-project/bckp-database/utils"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RouteBckpDatabase(app *fiber.App) {
	app.Get("/generate-token", GeneratToken)
	databaseGroup := app.Group("/bckp-database", ValidateBearerToken)
	databaseGroup.Get("/", GetAllLatestBackupDatabase)
	databaseGroup.Get("/:db_name", GetHistoryBackupByName)
	databaseGroup.Post("/:db_name", UploadFileHandler)
	databaseGroup.Get("/:id_file/download", DownloadFileHandler)
}

func validateZipFile(file *multipart.FileHeader) error {
	fileName := file.Filename

	if !strings.HasSuffix(fileName, ".zip") {
		return errors.New("invalid file format, only zip files allowed")
	}

	return nil
}

func UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("zip_file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			map[string]any{
				"message": err.Error(),
			},
		)
	}

	dbName := c.Params("db_name")

	fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

	err = validateZipFile(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(map[string]any{
				"message": err.Error(),
			})
	}

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

	c.Set("Content-Disposition", "attachment; filename=example.zip")
	c.Set("Content-Type", "application/zip")

	fmt.Println(zipFilePath)

	err := c.SendFile(zipFilePath)
	if err != nil {

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

func GetHistoryBackupByName(c *fiber.Ctx) error {

	database_name := c.Params("db_name")

	historiesDB, errCreateBckpDatabse := utils.BackupHistoryByName(database_name)

	if errCreateBckpDatabse != nil {
		logrus.Printf("Terjadi error : %s\n", errCreateBckpDatabse.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(map[string]any{
				"message": fmt.Sprintf("Server Error = %s", errCreateBckpDatabse.Error()),
			})
	}

	if len(historiesDB) == 0 {
		return c.Status(fiber.StatusNotFound).
			JSON(map[string]any{
				"message": fmt.Sprintf("Database  %s Not found! or no data history ", database_name),
			})
	}

	var result []map[string]interface{}
	for _, historie := range historiesDB {
		result = append(result, map[string]interface{}{
			"id":        historie.ID,
			"file_name": historie.FileName,
			"timestamp": historie.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"data": map[string]any{
			"database_name": database_name,
			"histories":     result,
		},
		"message": "Success",
	})

}

func GeneratToken(c *fiber.Ctx) error {

	token, err := utils.AddToken(model.Token{})

	if err != nil {
		logrus.Printf("Terjadi error : %s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(map[string]any{
				"message": fmt.Sprintf("Server Error = %s", err.Error()),
			})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"token":   token.Value,
		"message": "Success",
	})

}

func ValidateBearerToken(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"message": "Missing Authorization header",
		})
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"message": "Invalid Authorization format",
		})
	}

	token := parts[1]

	if len(token) != 36 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"message": "Invalid token length",
		})
	}

	uuiToken, err := uuid.Parse(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"message": "Invalid format token",
		})
	}

	_, err = utils.GetValueToken(uuiToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"message": "Token not found",
		})
	}

	return c.Next()
}
