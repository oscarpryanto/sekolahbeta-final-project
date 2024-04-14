package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Value     uuid.UUID `gorm:"unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (cr *Token) CreateToken(db *gorm.DB) error {
	err := db.
		Model(Token{}).
		Create(&cr).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (cr *Token) GetToken(db *gorm.DB) (Token, error) {
	res := Token{}
	err := db.Model(Token{}).Where("value = ?", cr.Value).First(&res).Error

	if err != nil {
		return Token{}, err
	}

	return res, nil
}
