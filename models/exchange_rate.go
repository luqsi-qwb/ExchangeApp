package models

import "time"

type ExchangeRate struct {
	ID           uint      `gorm:"primarykey" json:"_id"`
	FromCurrency string    `json:"fromCurrency" binding:"required"`
	ToCurrency   string    `json:"toCurrency" binding:"required"`
	Rate         float32   `json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
}
