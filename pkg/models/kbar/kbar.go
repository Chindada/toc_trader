// Package kbar package kbar
package kbar

import "gorm.io/gorm"

// Kbar Kbar
type Kbar struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	StockNum   string  `gorm:"column:stock_num;index:idx_kbar"`
	TimeStamp  int64   `gorm:"column:timestamp;index:idx_kbar"`
	Close      float64 `gorm:"column:close"`
	Open       float64 `gorm:"column:open"`
	High       float64 `gorm:"column:high"`
	Low        float64 `gorm:"column:low"`
	Volume     int64   `gorm:"column:volume"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Kbar) TableName() string {
	return "tick_kbar"
}
