package utils

import (
	"final-project/bckp-database/config"
	"final-project/bckp-database/model"
	"time"

	"github.com/google/uuid"
)

func AddToken(data model.Token) (model.Token, error) {
	token := uuid.New()
	data.Value = token
	data.CreatedAt = time.Now()
	err := data.CreateToken(config.Mysql.DB)

	return data, err
}

func GetValueToken(value uuid.UUID) (model.Token, error) {
	token := model.Token{
		Value: value,
	}
	return token.GetToken(config.Mysql.DB)
}
