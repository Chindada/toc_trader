// Package simulationcond package simulationcond
package simulationcond

import (
	"gorm.io/gorm"
)

// AnalyzeCondition AnalyzeCondition
type AnalyzeCondition struct {
	gorm.Model            `json:"-" swaggerignore:"true"`
	TrimHistoryCloseCount bool    `gorm:"column:trim_history_close_count" json:"trim_history_close_count"`
	HistoryCloseCount     int64   `gorm:"column:history_close_count" json:"history_close_count"`
	ForwardOutInRatio     float64 `gorm:"column:forward_out_in_ratio" json:"forward_out_in_ratio"`
	ReverseOutInRatio     float64 `gorm:"column:reverse_out_in_ratio" json:"reverse_out_in_ratio"`
	CloseChangeRatioLow   float64 `gorm:"column:close_change_ratio_low" json:"close_change_ratio_low"`
	CloseChangeRatioHigh  float64 `gorm:"column:close_change_ratio_high" json:"close_change_ratio_high"`
	OpenChangeRatio       float64 `gorm:"column:open_change_ratio" json:"open_change_ratio"`
	RsiHigh               float64 `gorm:"column:rsi_high" json:"rsi_high"`
	RsiLow                float64 `gorm:"column:rsi_low" json:"rsi_low"`
	TicksPeriodThreshold  float64 `gorm:"column:ticks_period_threshold" json:"ticks_period_threshold"`
	TicksPeriodLimit      float64 `gorm:"column:ticks_period_limit" json:"ticks_period_limit"`
	TicksPeriodCount      int     `gorm:"column:ticks_period_count" json:"ticks_period_count"`
	VolumePerSecond       int64   `gorm:"column:volume_per_second" json:"volume_per_second"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (AnalyzeCondition) TableName() string {
	return "simulate_cond"
}
