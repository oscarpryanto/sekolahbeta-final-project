package presenter

import (
	"final-project/bckp-database/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BckpDatabase struct {
	ID           uint      `json:"id"`
	DatabaseName string    `json:"database_name"`
	FileName     string    `json:"file_name"`
	TimeStamp    time.Time `json:"time_stamp"`
}

func AddBckpDatabseSuccessResponse(data *model.BckpDatabase) *fiber.Map {
	bckpDatabase := BckpDatabase{
		ID:           data.ID,
		DatabaseName: data.DatabaseName,
		FileName:     data.FileName,
		TimeStamp:    data.UpdatedAt,
	}
	return &fiber.Map{
		"data":    bckpDatabase,
		"message": "Success Import",
	}
}
